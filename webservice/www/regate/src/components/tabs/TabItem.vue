<template>
    <keep-alive v-if="!disabled" style="height: 100%;">
      <component 
        ref="componentMedia"
        :is="panelComposant"
        :sessionNumber="sessionNumber"
        :protocolObject="protocolObject"
        :id="panelid"
        style="height: 100%;"
      />
    </keep-alive>
  </template>
  
  <script>

  import NotFound from "./TabNotFound";
  import TabCanvas from "./TabCanvas";
  import TabTerminal from "./TabTerminal";
  import TabWelcome from "./TabWelcome";
  import TabDirectly from "./TabDirectly";
  import TabSettings from "./TabSettings";
  import { markRaw } from 'vue'


  export default {
    name: "TabItem",

    props: {
      disabled: {
        Type: Boolean,
        required: true,
        default: false,
      },
      sessionNumber: {
        Type: Number,
        required: true,
        default: 0,
      },
      protocolObject: {
        Type: Object,
        required: true,
        default: 0,
      },
      typePanel: {
        Type: String,
        required: true,
        default: "",
      }
    },
    data() {
      return {
        panelid:"p_"+this.sessionNumber
      };
    },
    mounted() {
    },
    computed: {
      panelComposant() {

        console.log("Create Panel session:",this.sessionNumber,"   typePanel:",this.typePanel,"Protocole:",this.protocolObject)
        switch(this.typePanel){
              case "Welcome":
                return markRaw(TabWelcome);

              // Not yet
              case "Directly":
                return markRaw(TabDirectly);
              case "Settings":
                return markRaw(TabSettings);
              default:
                break;
        }
        if(!this.protocolObject){
          console.log("typePanel not found and protocolObject not found");
            return markRaw(NotFound);
        }
        switch (this.protocolObject.typeUI) {
          default:
            console.log("Type panel not found:", this.protocolObject.typeUI, this);
            return markRaw(NotFound);
          case "canvas":
            console.log("Type Canvas:", this.protocolObject.typeUI, this);
            return markRaw(TabCanvas);
          case "terminal":
            console.log("Type terminal:", this.protocolObject.typeUI, this);
            return markRaw(TabTerminal);
        }
      },
    },
    methods: {
      keydown(e){
        console.log("panel keydown",this.protocolObject,e);
        if(this.protocolObject){
          this.protocolObject.keydown(e);
        }
      },
      keyup(e){
        console.log("panel keyup",this.protocolObject,e);
        if(this.protocolObject){
          this.protocolObject.keyup(e);
        }
      }
    }
  };
  </script>