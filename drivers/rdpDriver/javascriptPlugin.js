
class RDP {
    constructor(){}
    get name(){
        return "rdpDriver"
    }
    get insertTab(){
        return `
        `;
    }
    get typeUI(){
        return "canvas"
    }
    get icon(){
        return "rdp"
    }

    set(eventBus,canvasPanelVue){
        this.eventBus=eventBus;
        this.canvasPanelVue=canvasPanelVue;
        this.onUpdateCanvas();
    }

    onUpdateCanvas(){
        //
        var sessionNumber=this.canvasPanelVue.sessionNumber
        var idcanvas=this.canvasPanelVue.id

        // Render (force width and height)
        var canvas=document.getElementById(idcanvas);
        canvas.height=canvas.clientHeight;
        canvas.width=canvas.clientWidth;
        var render = new Mstsc.Canvas.create(canvas); 

        // adding eventBus listener for update canvas
        this.eventBus.$on(sessionNumber+"_Update", (messageWS) => {
            this.canvasPanelVue.remoteLoadingShow=false;
            var arrayBMP=messageWS.Msg
            for (var i in arrayBMP)
            { 
                var bitmap = arrayBMP[i];
                
                var raw = window.atob(arrayBMP[i].data);
                var rawLength = raw.length;
                var array = new Uint8Array(new ArrayBuffer(rawLength));

                for(i = 0; i < rawLength; i++) {
                    array[i] = raw.charCodeAt(i);
                }
                bitmap["data"]=array;
                render.update(bitmap);
            }
        })
    }
    keydown(e){
        console.log("RDP: keydown",e);
        
        this.canvasPanelVue.$ws.sendMessage({
            Session: this.canvasPanelVue.sessionNumber,
            Msg:{
                Type: "scancode",
                ScanCode:{
                    Button:Mstsc.scancode(e),
                    IsPressed:true
                }
        }});

        //e.preventDefault();
        e.stopPropagation() 
        return false;
    }
    keyup(e){
        console.log("RDP: keyup",e);
        this.canvasPanelVue.$ws.sendMessage({
            Session: this.canvasPanelVue.sessionNumber,
            Msg:{
                Type: "scancode",
                ScanCode:{
                    Button:Mstsc.scancode(e),
                    IsPressed:false
                }
        }});
        //e.preventDefault();
        e.stopPropagation();
        return false;
    }

    /**
     * Mouse button mapping
     * @param button {integer} client button number
     */
    mouseButtonMap  (button) {
        switch(button) {
            case 0:
                return 1;
            case 2:
                return 2;
            default:
                return 0;
        }
    }

    
    mousedownHandleEvent(e) {
        console.log("mousedown",this.sessionNumber,e);

        var coordMouseCanvas=this.canvasPanelVue.computeCoordMouseCanvas(e);
        this.canvasPanelVue.$ws.sendMessage({
            Session: this.canvasPanelVue.sessionNumber,
            Msg:{
                Type: "mouse",
                ScanCode:{
                    X:coordMouseCanvas.X,
                    Y:coordMouseCanvas.Y,
                    Button: this.mouseButtonMap(e.button),
                    IsPressed:true
                }
        }});
            
        e.preventDefault();
        return false;
    }
    mouseupHandleEvent(e) {
        //if (!self.activeSession) return;
        console.log("mouseup",this.canvasPanelVue.sessionNumber,e);
            
        var coordMouseCanvas=this.canvasPanelVue.computeCoordMouseCanvas(e);
        this.canvasPanelVue.$ws.sendMessage({
            Session: this.canvasPanelVue.sessionNumber,
            Msg:{
                Type: "mouse",
                ScanCode:{
                    X:coordMouseCanvas.X,
                    Y:coordMouseCanvas.Y,
                    Button:this.mouseButtonMap(e.button),
                    IsPressed:false
                }
        }});
        
        e.preventDefault();
        return false;
    }
    mousemoveHandleEvent (e){
        console.log("mouse event",this.canvasPanelVue.sessionNumber,e);

        var coordMouseScreen=this.canvasPanelVue.computeCoordMouseCanvas(e);
        console.log(coordMouseScreen);
        this.canvasPanelVue.$ws.sendMessage({
            Session: this.canvasPanelVue.sessionNumber,
            Msg:{
                Type: "mouse",
                ScanCode:{
                    X:coordMouseScreen.X,
                    Y:coordMouseScreen.Y,

                    Button:0,
                    IsPressed:false
                }
        }});
    }
    showContextMenu (e) {
        console.log("contextMenu",this.canvasPanelVue.sessionNumber,e);

        var offset = this.canvasPanelVue.elementOffset();
        this.canvasPanelVue.$ws.sendMessage({
            Session: this.canvasPanelVue.sessionNumber,
            Msg:{
                Type: "mouse",
                ScanCode:{
                    X:Math.floor((e.clientX - offset.left)*1.018),
                    Y:Math.floor((e.clientY - offset.top)*1.02),
                    Button:this.mouseButtonMap(e.button),
                    IsPressed:false
                }
        }});
        e.preventDefault();
        return false;
    }
    handleWheel (e) {
        console.log("contextMenu",e);

        var isHorizontal = false;
        var delta = e.detail;
        var step = Math.round(Math.abs(delta) * 15 / 8);
        
        var offset = this.canvasPanelVue.elementOffset();
        this.canvasPanelVue.$ws.sendMessage({
            Session: this.canvasPanelVue.sessionNumber,
            Msg:{
                Type: "mouse",
                Wheel:{
                    X:Math.floor((e.clientX - offset.left)*1.018),
                    Y:Math.floor((e.clientY - offset.top)*1.02),
                    Step:step,
                    IsNeg:delta>0,
                    IsH:isHorizontal
                }
        }});
        e.preventDefault();
        return false;
    }

    // fonction option configuration return string: Error or object for configuration connection WS
    static encodeConfiguration(dataForm){

        if(dataForm.user=="" || dataForm.user==undefined) {
            return ("user is empty")
        }
        if(dataForm.domain==""|| dataForm.domain==undefined) {
            return ("domain is empty")
        }
        if(dataForm.host==""|| dataForm.host==undefined) {
            return ("host is empty")
        }

        var userdomain=dataForm.domain+"/"+dataForm.user
        var encodedSchema = encodeURIComponent("rdp");
        var encodedAuthority = `${encodeURIComponent(userdomain)}:${encodeURIComponent(dataForm.password)}@${encodeURIComponent(dataForm.host)}`;
//        const encodedPath = encodeURIComponent(path);
//        const encodedQuery = encodeURIComponent(query);
//        var encodedFragment = encodeURIComponent(fragment);

//        return `${encodedSchema}://${encodedAuthority}${encodedPath}?${encodedQuery}#${encodedFragment}`;
        return {
            Name: dataForm.name,
            URL:`${encodedSchema}://${encodedAuthority}`    
        }
    }
    
}


listPlugin["rdp"]=RDP;

