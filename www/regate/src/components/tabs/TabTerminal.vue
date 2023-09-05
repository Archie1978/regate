<template>
    <v-row :style="{'overflow':'scroll', 'width': remoteLoadingWidth + 'px','height': remoteLoadingHeight + 'px' }" class="tab">
        <div v-bind:id="id" ref="terminal" style="height: 100%;flex:auto;"></div>
        <div class="remoteLoading" 
            @contextmenu = "showContextMenu($event)"
            ref="loading"
            v-show="remoteLoadingShow"
            :style="{'backgroundImage':  'url('+remoteLoadingImg+')', 'top': remoteLoadingTop + 'px','left': remoteLoadingLeft + 'px','width': remoteLoadingWidth + 'px','height': remoteLoadingHeight + 'px' }">
        </div>
    </v-row>
</template>

<script>

import eventBus from '../../eventBus'
import loadlogo from '@/assets/loading.gif'

import '../../../node_modules/xterm/css/xterm.css'
import { Terminal } from 'xterm';

import { WebLinksAddon } from 'xterm-addon-web-links';

export default {
    name: "TabTerminal",
    components: {},
    data: () => ({
        valid: true,
        name: "Terminal"  ,
        
        remoteLoadingShow: true,
        remoteLoadingImg:loadlogo,
        remoteLoadingTop:"0px",
        remoteLoadingLeft:"0px",
        remoteLoadingHeight:"0px",
        remoteLoadingWidth:"0px",
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

            logLevel: 1, // 1:debug  5:Off
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
        showContextMenu: function (e) {
            console.log("contextMenu",this.sessionNumber,e);
            if(this.protocolObject){
                if(this.protocolObject.showContextMenu){
                    this.protocolObject.showContextMenu();
                }
            }
        }
    }

};

</script>
