<template>

    <v-row :style="{'overflow':'scroll', 'width': remoteLoadingWidth + 'px','height': remoteLoadingHeight + 'px' }" class="tab">
        <div v-bind:id="id"
            ref="terminal"
            style="height: 100%;flex:auto;" 
            @contextmenu = "showContextMenu($event)"
            >
        </div>
        <div class="remoteLoading" 
            ref="loading"
            v-show="remoteLoadingShow"
            :style="{'backgroundImage':  'url('+remoteLoadingImg+')', 'top': remoteLoadingTop + 'px','left': remoteLoadingLeft + 'px','width': remoteLoadingWidth + 'px','height': remoteLoadingHeight + 'px' }">
        </div>
        <div ref="menucontexuel"  v-show="showMenu" 
            @mouseleave="showMenu=false"
            :style="{ 'top': menuY + 'px','left': menuX + 'px','width': 300 + 'px','height': 100 + 'px', position:'absolute',background: 'white','z-index':'100' }">
            <template  v-for="item in menucontexuelItems" :key="item.name">
                <button @click="menucontexuelClick(  item )" style="width:100%" > {{ item.name }}</button>
            </template>
        </div>
    </v-row>
    <div ref="info"></div>
</template>

<style>

button:hover {
  background-color: gold;
}
</style>

<script>

import eventBus from '../../eventBus'
import loadlogo from '@/assets/loading.gif'

import '../../../node_modules/xterm/css/xterm.css'
import { Terminal } from 'xterm';

import { WebLinksAddon } from 'xterm-addon-web-links';

//import VueSimpleContextMenu from 'vue-simple-context-menu';
//import 'vue-simple-context-menu/dist/vue-simple-context-menu.css';



export default {
    name: "TabTerminal",
    components: {
        //VueSimpleContextMenu
    },
    data: () => ({
        valid: true,
        name: "Terminal"  ,
        
        remoteLoadingShow: true,
        remoteLoadingImg:loadlogo,
        remoteLoadingTop:"0px",
        remoteLoadingLeft:"0px",
        remoteLoadingHeight:"0px",
        remoteLoadingWidth:"0px",

        showMenu: false,


        // Menu into terminal
        menuX:0,
        menuY:0,
        menucontexuelItems: [],

    }),
    props: {
        sessionNumber: {
        Type: String,
        required: true,
        default: "",
      },
      protocolObject: {
        Type: Object,
        required: true,
        default: {},
      },
      id: {
        Type: String,
        required: true,
        default: "terminal_0",
      },
    },

    
    mounted() {

        var tabpanel=document.getElementById("tabpanel");
        var cols=Math.floor(tabpanel.offsetWidth/9.2);//200,
        var rows=Math.floor(tabpanel.offsetHeight/18.5);//40,

        let t=new Terminal({
            role: "server",
            shell: "linux",

            cols: cols,
            rows: rows,

            convertEol: true,
            cursorBlink: true,

            logLevel: 4, // 1:debug  5:Off
            //visualBell: true,
        });

        // Set protocolObject
        this.protocolObject.set(eventBus,this,t);

        // Start connection with option tty
        this.$ws.sendMessage({
            Session: this.sessionNumber,
            Msg:{
                Type: "size",
                Screen:{
                    Cols:cols,
                    Rows:rows,
                }
        }});

        // Addon
        t.loadAddon(new WebLinksAddon());

        // Add loading
        this.remoteLoadingTop=tabpanel.offsetTop;
        this.remoteLoadingLeft=tabpanel.offsetLeft;
        this.remoteLoadingWidth=tabpanel.offsetWidth;
        this.remoteLoadingHeight=tabpanel.offsetHeight;

    
    },



    methods: {
        // Add context menu
        showContextMenu: function (event) {
            // Get Menu
            if(this.protocolObject){
                if(this.protocolObject.showContextMenu){
                    this.menucontexuelItems=this.protocolObject.showContextMenu(event);
                }
            }

            if(this.menucontexuelItems.length>0){
                // Display and move menu
                this.showMenu = true;
                this.menuX = event.x>10?event.x-10:event.x;
                this.menuY = event.y>10?event.y-10:event.y;

                // There is a menu
                event.stopPropagation();
                event.preventDefault();
            }
        },
        menucontexuelClick(item){
            console.log("menucontexuelClick",item);
            this.getSelectionIntoTerminal();
            if(this.protocolObject.showContextMenu){
                this.protocolObject.menucontexuelClick(item);
            }
            this.showMenu = false;
        },

        // Get all selection into terminal
        getSelectionIntoTerminal(){
            let t="";
            let childrenWithClass = this.$refs.terminal.querySelectorAll(".xterm-decoration-top");

            childrenWithClass.forEach((element) => {
                // Faites quelque chose avec chaque élément, par exemple, imprimez le texte
                t+=element.textContent;
            });
            return t;
        }
    }

};

</script>
