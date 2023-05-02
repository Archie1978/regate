class Ssh {
    constructor(){
        this.connected=false;
        this.terminal=null;
        this.terminalVue=null;
    }
    get name(){
        return "sshDriver"
    }
    get insertTab(){
        return `
        
        `;
    }
    get typeUI(){
        return "terminal"
    }
    get icon(){
        return "ssh"
    }

    set(eventBus,terminalVue,terminal){
        this.terminal=terminal;
        this.terminalVue=terminalVue;
        this.eventBus=eventBus;
        this.terminalElement=document.getElementById(terminalVue.id);

        // Create terminal
        this.terminal.open(this.terminalElement);

        var sshObject=this;
        this.terminal.onKey((e) => {
            console.log("this.terminal.onKey",e);

            if(this.connected){
                const ev = e.domEvent;
                const printable = !ev.altKey && !ev.ctrlKey && !ev.metaKey;

                console.log("terminal",e);
                var c=e.key;
                var charprint=e.key;
                switch(e.key){
                    case "Enter":
                        c=String.fromCharCode(10);
                        charprint="\n";
                        break;
                    case "Backspace":
                        c=String.fromCharCode(8);
                        break;
                    case "Tab":
                        c=String.fromCharCode(9);
                        break;
                    case "ArrowUp":
                        c=String.fromCharCode(38);
                        break;
                    case "ArrowDown":
                        c=String.fromCharCode(39);
                    case "Space":
                        c=" ";
                        break;
                    defaut:
                        break;
                }

                this.terminalVue.$ws.sendMessage({    
                    Session: this.terminalVue.sessionNumber,
                    Msg:{
                        Type: "key",
                        Key:{
                          Keys:[c],
                        }
                }});
            }
            return false;
        });
        this.focus();

        // Set stdout event
        this.stdout(this.eventBus,terminalVue.sessionNumber);
    }

    base64ToArrayBuffer(base64) {
        var binary_string = window.atob(base64);
        var len = binary_string.length;
        var bytes = new Uint8Array(len);
        for (var i = 0; i < len; i++) {
            bytes[i] = binary_string.charCodeAt(i);
        }
        return bytes.buffer;
    }

    stdout(eventBus,sessionNumber){
        // Send from webservice
        eventBus.$on(sessionNumber+"_Out", (messageWS) => {
            console.log("out",messageWS)

            // Remove masque
            if(!this.connected){
                console.log("Hide loading")
                this.connected=true;
                this.terminalVue.remoteLoadingShow=false;
            }
            
            const sdata = atob(messageWS.Msg.Content);
            const bytes = new Uint8Array(sdata.length);
            for (let i = 0; i < bytes.length; ++i) {
                bytes[i] = sdata.charCodeAt(i);
            }
              
            // write bytes to terminal
            this.terminal.write(bytes);

            //this.terminal.write(atob(messageWS.Msg.Content));
        });
    }

    keydown(e){
        console.log("TerminalSsh: keydown",this.terminal,e,document.getElementById(this.terminalVue.id));
        return false;
    }
    
    keyup(e){
        return false;
    }

    focus(){
        if(this.terminalVue) {
            document.getElementById(this.terminalVue.id).querySelector("textarea").focus();
        }
    }
}


listPlugin["ssh"]=Ssh;

