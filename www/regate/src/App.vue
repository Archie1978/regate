<template>
  <VApp>
    <SlideBar></SlideBar>
    <VMain>
      <TabPanels/>
    </VMain>

    <ModalConfigurationConnection v-show="showModalConfigurationConnection"  @close="modalConfigurationConnectionClose"/>
    <ModalMessage v-show="showModalMessage" :title="modalMessageTitle" :message="modalMessageBody" @close="modalMessageClose"/>

    <div v-if="lostConnection" style="justify-content: center;align-items: center;background-color: gray;z-index: 1500;position:absolute;display:flex;width: 100%;height: 100%;  margin: auto;">
      <div>
        LOST CONNECTION
      </div>
    </div>

  </VApp>

</template>


<style>
.panelmain{
  width:100%;
  height: 100%;

  border: black;
  border-top-style: none;
  border-right-style: none;
  border-bottom-style: none;
  border-left-style: none;
  border-style: groove;
}

.remoteLoading{
  position: absolute;
  background-image: "@/assets/loading.png";
  width: 1000px;
  background-color: white;
  background-position: center;
  background-repeat: no-repeat;
  opacity: 0.6;
  z-index: 100;
}

.tab {
  margin-top: -12px!important;
}

</style>




<script>
import SlideBar from "./components/SlideBar.vue"
import TabPanels from "./components/TabPanels.vue"
import eventBus from './eventBus'
import ModalConfigurationConnection from './components/ModalConfigurationConnection.vue'
import ModalMessage from './components/ModalMessage.vue'


export default {
  name: 'App',

  components: {
    SlideBar,
    TabPanels,
    ModalConfigurationConnection,
    ModalMessage,
  },

  data: () => ({
    showModalConfigurationConnection: false,

    showModalMessage: false,
    modalMessageTitle: "",
    modalMessageBody: "",
    
    lostConnection: false,
  }),

  methods:{
    modalConfigurationConnectionClose(){
      this.showModalConfigurationConnection=false;
    },
    modalMessageClose(){
      this.showModalMessage=false;
    }
  },

  mounted:function(){

    // Add script from plugin
    const scriptlocale = document.createElement("script");
    scriptlocale.setAttribute(
      "src",
      "http://localhost:8088/addon-local.js"
    )
    document.head.appendChild(scriptlocale);
    scriptlocale.onload = () => {
        console.log("script addon-local.js loaded");
    }

  //

    // Get Signal event into bus
    var app=this;
    eventBus.$on("WebSocketLost", () => {
      app.lostConnection=true;
    })

    eventBus.$on("ShowModalConfigurationConnection", () => {
      this.showModalConfigurationConnection=true;
    })

    eventBus.$on("ShowModalMessage", (title,body) => {
      console.log(title,body)
      this.modalMessageTitle=title;
      this.modalMessageBody=body;
      this.showModalMessage=true;
    })


    document.body.addEventListener('copy', function(e){
      console.log(e);
    });

  }
}
</script>
