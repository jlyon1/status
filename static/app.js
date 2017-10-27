var textbox = Vue.component('text-box',{
  template:`<div class="message_body">
  <div> <input v-model="message" placeholder="Enter a msg"/> </div>
  <div style="margin-top: 10px;"> <input type="password" v-model="password" placeholder=""/> </div>
  <div style="margin-top: 10px; float:right;"> <button v-on:click="send">Submit</button> </div>
  </div>`,
  data (){
    return {
      message: "Hello World",
      password: "",
      counter: 0,
      ws: {}
    }
  },
  methods: {
    send: function(){
      console.log("asdf");
      var val = {
        message: this.message, password: this.password
      }
      $.post("update/submit/",JSON.stringify(val),function(data){
        console.log(data)
      })
    }
  }
})
var msgbdy = Vue.component('message-body',{
  template:`<div class="message_body">{{msgtxt}}</div>`,
  data (){
    return {
      msgtxt: "msg",
      ws: {}
    }
  },
  methods:{
    openWebsocket: function(){
      if (ws) {
        return false;
      }
      var ws = new WebSocket("ws://127.0.0.1:8080/ws");

      ws.onopen = function(evt) {
        console.log("ya did it")
      }
      ws.onclose = function(evt) {
        ws = null;
      }
      var el = this;
      ws.onmessage = function(evt) {
        console.log(this.msgtxt);
        el.msgtxt = evt.data;
      }
      ws.onerror = function(evt) {
        console.log(evt.data)
      }

    }
  },
  mounted (){
    this.openWebsocket();
  }

})

var titlebar = Vue.component('title-bar',{
  props: ['msg'],
  template:`<div class="titlebar">{{msg}}</div>`,
  data (){
    return {
    }
  }

})

var ShuttleTracker = new Vue({
  el: '#app-vue',
  data: {
  }

});
