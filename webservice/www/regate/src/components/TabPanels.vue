<template>
  <div class="contentTab">
    <v-row>
      <v-col cols="13">
        <vue-tabs-chrome
          v-model="tab"
          :tabs="tabs"
          ref="tabs"
          :minHiddenWidth="120"
          :theme="theme"
          @remove="eventTabRemove"
        />
      </v-col>
    </v-row>
    <v-row style="height: 90%;margin: 3px;">
      <v-col style="height: 100%;" id="tabpanel">
        <Tab
          v-for="tabPanel in tabPanels"
          :ref="tabPanel.key"
          :key="tabPanel.key"
          v-show="(tagTarget == tabPanel.key)"
          :typePanel="tabPanel.typePanel"
          :sessionNumber="tabPanel.sessionNumber"
          :protocolObject="tabPanel.protocolObject"
          style="height: 100%;"
        />
      </v-col>
    </v-row>
  </div>
</template>

<style>
.contentTab {
  height: 100%;
}
</style>
  
<script>
import Vue3TabsChrome from 'vue3-tabs-chrome'
import 'vue3-tabs-chrome/dist/vue3-tabs-chrome.css'
  
import Tab from "../components/tabs/TabItem";

export default {
  name: "TabPanels",
  components: {
    "vue-tabs-chrome": Vue3TabsChrome,
    Tab: Tab,
  },
  props: {},
  methods: {

    // Add remote tab
    addTabRemote(sessionNumber,label,indexkey,protocolObject) {
      if(this.existTab(indexkey)){
        return
      }

      //  Add tab
      let newTabs = [
        {
          label: label,
          key: indexkey,
          favicon3: protocolObject.icon,
          favicon: protocolObject.icon,
        },
      ];

      // Add content tab
      this.$refs.tabs.addTab(...newTabs);
      this.tabPanels.push({
        sessionNumber: sessionNumber,
        protocolObject: protocolObject,
        key: indexkey,
        enable: false,
      });

      // Active new tab
      this.tab = indexkey;

    },

    // Add tab setting ....
    addTabCustom(label,typePanel,icon) {
      if(this.existTab(label)){
        return
      }

      this.$refs.tabs.addTab({
          label: label,
          key: label,
          favicon3: icon,
          favicon: icon,
        }
      );
      this.tabPanels.push({
        sessionNumber: 0,
        key: label,
        typePanel:typePanel,
        enable: false,
      }
      );


      // Active new tab
      this.tab = label;

    },

    // Tab Exist
    existTab(indexkey){
      for (var i = 0; i < this.tabPanels.length; i++) {
        if(this.tabPanels[i].indexkey==indexkey){
          return true;
        }
      }
      return false;
    },


    // Get CurrentTab
    getCurrentTab() {
      return this.tabs.find((item) => item.key === this.tab);
    },
    getCurrentTabPanel() {
      return this.tabPanels.find((item) => item.key === this.tab);
    },

    // Remove Tab by tab
    eventTabRemove(tab) {
      console.log("eventTabRemove",tab);
      this.removeTab(tab.key);

      // Send webservice session End
      this.$ws.sendMessage({    
            Session: this.sessionNumber,
            Msg:{
                Type: "key",
                Command:"Close"
      }});
    },


    // Remove tab by index Key
    removeTab(indexkey){
      console.log("removeTab",this,indexkey);

      // Remove panel intoo tab
      for (var i = 0; i < this.tabPanels.length; i++) {
        if (this.tabPanels[i].key === indexkey) {

          // selet new tab title
          if(i==0){
            if(i+1<this.tabPanels.length){
              this.tab=this.tabPanels[i+1].key;
            }else{
              this.tab="";
            }
          }else{
            this.tab=this.tabPanels[i-1].key;
          }
          this.tabPanels.splice(i, 1);

          // Remove title tab
          this.$refs.tabs.removeTab(indexkey);

          return;
        }
      }
    },
  },
  computed: {
    theme() {
      //return this.$store.state.userSettings.darkMode ? "dark" : "";
      return "dark";
    },
  },
  watch: {
    tab() {

      console.log("tab:",this,this.tab);
      let tab = this.getCurrentTab();
      if(!tab){
        return "";
      }
      let tagTarget = tab.key || "";
      if (!tab.key) {
        return "";
      }
      this.tagTarget = tagTarget;

      // Active panel
      let panelcurrent=this.getCurrentTabPanel();
      if(panelcurrent.protocolObject){
        if(panelcurrent.protocolObject.focus){
          panelcurrent.protocolObject.focus();
        }
      }
    },
  },
  data() {
    return {
      tab: "",
      tagTarget: "",
      tabs: [],
      tabPanels: [],
    };
  },
  mounted() {
    this.$ws.tabPanel=this;

    // bind keyboard event to tab
    var tabpanel=this;
    window.addEventListener('keydown', function (e) {
        var tabref = tabpanel.$refs[tabpanel.tagTarget];
        if(tabref){
          if(tabref.length>0) {
            return tabref[0].keydown(e);
          }
        }
        //e.preventDefault();
        return false;
    });
    
    window.addEventListener('keyup', function (e) {
        console.log("keyup",e);

        var tabref = tabpanel.$refs[tabpanel.tagTarget];
        if(tabref){
          if(tabref.length>0) {
            return tabref[0].keyup(e);
          }
        }
        //e.preventDefault();
        //return false;
    });
  },
};
</script>