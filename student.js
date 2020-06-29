new Vue({
    el: '#app',

    data: {
        ws: null, // Our websocket
        wrd: '', // Holds words to be sent to the server
        rollno: null, // roll number
        about: null, // student writing area
        joined: false // True if rollno and about have been filled in
    },

    created: function() {
        var self = this;
        this.ws = new WebSocket('localhos:8080/ws');
        this.ws.addEventListener('message', function(e) {
            var msg = JSON.parse(e.data);
            self.Content += '<div class="chip">'
                    + msg.about
                + '</div>'
                
            var element = document.getElementById('Type Area');
            element.scrollTop = element.scrollHeight; 
        });
    },

    methods: {
        send: function () {
            if (this.newMsg != '') {
                this.ws.send(
                    JSON.stringify({
                        rollno: this.rollno,
                        about: this.about,
                        message: $('<p>').html(this.newMsg).text() 
                    }
                ));
                this.newMsg = ''; 
            }
        },
        
        join: function () {
            if (!this.rollno) {
                Materialize.toast('You must enter an roll number', 2000);
                return
            }
            if (!this.about) {
                Materialize.toast('start typing about yourself', 2000);
                return
            }
            this.rollno = $('<p>').html(this.rollno).text();
            this.about = $('<p>').html(this.about).text();
            this.joined = true;
        },
        
    }
});