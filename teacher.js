
new Vue({
    el:'#get-data-from-server-app',
	data: {
        ws: null,
        locat: null,
        error: null
        },

        created: function(){
            var self = this;
            this.ws = new WebSocket('localhos:8080/ws');
		    this.init();
        },
        
        methods: {
		init: function(){
		this.loadData();
		},
        
        loadData: function() {
            this.$http.get().then((response) => {
                if(!!response.body) {
				this.locat = response.body;
				}
						}, (response) => {
							this.error = response;
						});
            
            
					}
				}
}); 