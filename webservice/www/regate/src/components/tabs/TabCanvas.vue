<template>
    <v-row v-show="!isLoading" style="margin-top: -12px;">
        <!--
        <canvas class="panelmain"
            ref="canvas"
            v-bind:id="id"
            v-on:mousedown='mousedownHandleEvent'
            v-on:mouseup='mouseupHandleEvent'
            v-on:mousemove="mousemoveHandleEvent"
            @contextmenu = "showContextMenu($event)"
            v-on:wheel="handleWheel"
            :style="{ position: 'absolute','top': remoteLoadingTop + 'px','left': remoteLoadingLeft + 'px','width': remoteLoadingWidth + 'px','height': remoteLoadingHeight + 'px' }">
        </canvas>
    -->
        <canvas class="panelmain"
            ref="canvas"
            v-bind:id="id"
            v-on:mousedown='mousedownHandleEvent'
            v-on:mouseup='mouseupHandleEvent'
            v-on:mousemove="mousemoveHandleEvent"
            @contextmenu = "showContextMenu($event)"
            v-on:wheel="handleWheel"
        >
        </canvas>
        <div class="remoteLoading" 
            ref="loading"
            v-show="remoteLoadingShow"
            :style="{'backgroundImage':  'url('+remoteLoadingImg+')', 'top': remoteLoadingTop + 'px','left': remoteLoadingLeft + 'px','width': remoteLoadingWidth + 'px','height': remoteLoadingHeight + 'px' }">
        </div>
    </v-row>
</template>

<script>

//:style="{ 'backgroundImage':  'url('+remoteLoadingImg+')','top': remoteLoadingTop + 'px','left': remoteLoadingLeft + 'px','width': remoteLoadingWidth + 'px','height': remoteLoadingHeight + 'px' }">

import eventBus from '../../eventBus'
import loadlogo from '@/assets/loading.gif'
export default {
    name: "TabCanvas",
    components: {},

    data: () => ({
        valid: true,
        name: "",
  
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
        default: "canvas_0",
      },
    },

    methods: {
        /**
		 * Compute screen offset for a target element
		 * @return {top : {integer}, left {integer}}
		 */
		elementOffset : function () {
            var el = this.$refs.canvas;
            var x = 0;
            var y = 0;
            while (el && !isNaN( el.offsetLeft ) && !isNaN( el.offsetTop )) {
                x += el.offsetLeft - el.scrollLeft;
                y += el.offsetTop - el.scrollTop;
                el = el.offsetParent;
            }
            console.log("elementOffset:  x:",x,"y:",y);
            return { top: y, left: x };
		},

        /*
         *  Get coordonne mouse into Canvas
         */
        computeCoordMouseCanvas:function(e){
            var offset = this.elementOffset();
            return {
                //X:Math.floor((e.clientX - offset.left)*1.018),
                //Y:Math.floor((e.clientY - offset.top)*1.02),
                X:Math.floor((e.clientX - offset.left)),
                Y:Math.floor((e.clientY - offset.top)),
            }
        },

        /**
         * Mouse button mapping
         * @param button {integer} client button number
         */
        mouseButtonMap : function (button) {
            switch(button) {
                case 0:
                    return 1;
                case 2:
                    return 2;
                default:
                    return 0;
            }
        },
        mousedownHandleEvent: function(e) {
            if(this.protocolObject){
                if(this.protocolObject.mousedownHandleEvent){
                    this.protocolObject.mousedownHandleEvent(e);
                }
            }
        },
        mouseupHandleEvent: function(e) {
            if(this.protocolObject){
                if(this.protocolObject.mouseupHandleEvent){
                    this.protocolObject.mouseupHandleEvent(e);
                }
            }
        },
        mousemoveHandleEvent: function (e){
            if(this.protocolObject){
                if(this.protocolObject.mousemoveHandleEvent){
                    this.protocolObject.mousemoveHandleEvent(e);
                }
            }
        },
        showContextMenu: function (e) {
            if(this.protocolObject){
                if(this.protocolObject.showContextMenu){
                    this.protocolObject.showContextMenu(e);
                }
            }
        },
        handleWheel:function (e) {
            if(this.protocolObject){
                if(this.protocolObject.handleWheel){
                    this.protocolObject.handleWheel(e);
                }
            }
        }
    },
    mounted() {
        // Fixe size of canvas
        var tabpanel=document.getElementById("tabpanel");
        var canvas=this.$refs.canvas;
        canvas.width=tabpanel.offsetWidth;
        canvas.height=tabpanel.offsetHeight;

        // Update protocolObject with canvas
        this.protocolObject.set(
            eventBus,
            this,
        );

        // Set size panel for connexion
        console.log("canvas",canvas);
        this.$ws.sendMessage({
            Session: this.sessionNumber,
            Msg:{
                Type: "size",
                Screen:{
                    Width:canvas.clientWidth,
                    Height:canvas.clientHeight,
                }
        }});

        // Add loading
        this.remoteLoadingTop=canvas.offsetTop;
        this.remoteLoadingLeft=canvas.offsetLeft;
        this.remoteLoadingWidth=canvas.offsetWidth;
        this.remoteLoadingHeight=canvas.offsetHeight;
    },
};


</script>
