class Ssh {
    constructor(){
        this.connected=false;
        this.terminal=null;
        this.terminalVue=null;
        this.pathupload="";
    }

    // Driver name
    get name(){
        return "sshDriver"
    }

    // Add tab
    get insertTab(){
        return ``;
    }

    // Start UI console terminal for connexion
    get typeUI(){
        return "terminal"
    }
    get icon(){
        return "ssh"
    }

    // connect Console with sh ssh
    set(eventBus,terminalVue,terminal){
        this.terminal=terminal;
        this.terminalVue=terminalVue;
        this.eventBus=eventBus;
        this.terminalElement=document.getElementById(terminalVue.id);

        // Create terminal
        this.terminal.open(this.terminalElement);

        // Create 
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
        this.stdout(this.eventBus,this.terminalVue.sessionNumber);


        // Add event drop file
        let rowElements = this.terminalVue.$refs.terminal.querySelectorAll("div");
        let me=this;
        rowElements.forEach(row => {
            row.addEventListener(
            "drop",
            (e) => {
                e.preventDefault();
                e.stopPropagation();

                if(me.pathupload==""){
                    alert("Set path uploas with menu contextuel")
                    return
                }

                const uploadForm = new FormData();
                for (const file of e.dataTransfer.files) {
                    uploadForm.append('file', file);
                }

                fetch('/uploadFile?sessionNumber='+me.terminalVue.sessionNumber+"&pathDir="+encodeURI(me.pathupload), {
                    method: 'POST',
                    body: uploadForm
                })
                .then(response => response.json())
                .then(data => {
                    console.log('Success:', data);
                })
                .catch((error) => {
                    console.error('Error:', error);
                });


            },
            false
        )
        });
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

    // Cursor function
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

    // Menucontextuel: return menu for configure drap and drop
    showContextMenu(e){
        console.log("TerminalSsh: contextMenu:",this.terminalVue,e);
        //document.getElementById(this.terminalVue.id).querySelector("textarea").focus();
        return         [
            {
                name: 'Download path',
                slug: 'add-star',
            },
            {
                name: 'Set Dir Upload',
                slug: 'remove-star',
            },
        ];
    }

    // Action menu contextuel
    menucontexuelClick(item){
        console.log("TerminalSsh: menucontexuelClick",item,this);
        let text = this.terminalVue.getSelectionIntoTerminal();
        let sessionNumber=this.terminalVue.sessionNumber;
        switch(item.name){
            case "Download path":
                this.downloadFile(encodeURI("/downloadFile?sessionNumber="+sessionNumber+"&path="+text),"download")
            case "Set Dir Upload":
                this.pathupload=text;
            default:
                console.error("RemoteSsh: menucontexuelClick: ",item.name)
        }
        console.log("TerminalSsh: menucontexuelClick:",window.getSelection(), text)
    }

    // fonction option configuration return string: Error or object for configuration connection WS
    static encodeConfiguration(dataForm){
        if(dataForm.user==""||dataForm.user==undefined) {
            return ("user is empty")
        }
        if(dataForm.host==""||dataForm.host==undefined) {
            return ("host is empty")
        }
        if(dataForm.port==""||dataForm.port==undefined) {
            return ("port is empty")
        }
        let encodedSchema = encodeURIComponent("ssh");
        let encodedAuthority = `${encodeURIComponent(dataForm.user)}:${encodeURIComponent(dataForm.password)}@${encodeURIComponent(dataForm.host)}:${dataForm.port}`;
    
        return {
            Name: dataForm.name,
            URL:`${encodedSchema}://${encodedAuthority}`
        }
    }


    downloadFile(url, nomFichier) {
        const downloadLink = document.createElement('a');
        downloadLink.href = url;
        downloadLink.download = nomFichier;
    
        document.body.appendChild(downloadLink);
        downloadLink.click();
        document.body.removeChild(downloadLink);
    }

    
    
}


listPlugin["ssh"]=Ssh;

