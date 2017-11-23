
var msgbdy = Vue.component('message-body',{
  template:`<div class="message_body">{{msgtxt}}</div>`,
  data (){
    return {
      msgtxt: "msg",
      ws: {}
    }
  },
  methods:{

  },
  mounted (){

  }

})

Vue.component("title-bar",{
  props: ['txt'],
  template: `<div v-bind:style=titleStyle><p v-bind:style=paragraphStyle>{{txt}}</p></div>`,
  data (){
    return{
      titleStyle: {position:"absolute",backgroundColor:"#eee",height:"34px",width:"auto",top:"1",left:"0",right:"0"},
      paragraphStyle: {float: "left",height:"34px",lineHeight:"34px",verticalAlign:"center",paddingLeft:"30px",margin:"0"},
      titleText: "Joseph Lyon"
    }
  }

});

var ShuttleTracker = new Vue({
  el: '#app-vue',
  data: {
  }

});
