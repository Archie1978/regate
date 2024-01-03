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

            // Send key
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

        // Add focus into terminal
        this.focus();

        // Set stdout event
        this.stdout(this.eventBus,this.terminalVue.sessionNumber);


        // Add event drop file
        let rowElements = this.terminalVue.$refs.terminal.querySelectorAll("div");
        let me=this;
        rowElements.forEach(row => {

            // Display red border
            row.addEventListener(
                "dragover",
                (e) => {
                    // Ajuste
                    if(!e.target.classList.contains("xterm-viewport")){
                        // Add red into xterm
                        if(e.target.parentElement.classList.contains("xterm-rows")) {
                            e.target.parentElement.style.borderColor = "red";
                            e.target.parentElement.style.borderStyle = "solid";
                        }else{
                            // It is a case xterm, up to xterm tag
                            e.target.parentElement.parentElement.style.borderColor = "red";
                            e.target.parentElement.parentElement.style.borderStyle = "solid";
                        }
                    }
                }
            );

            // Discard red border
            row.addEventListener(
                "dragleave",
                (e) => {
                    // Delete border without
                    e.target.parentElement.style.borderStyle = "none";
                    e.target.parentElement.parentElement.style.borderStyle = "none";
                }
            );

            // Send file
            row.addEventListener(
                "drop",
                (e) => {
                    // remove event propagation
                    e.preventDefault();
                    e.stopPropagation();

                    // Remove style
                    e.target.parentElement.parentElement.style.borderStyle = "none";
                    e.target.parentElement.style.borderStyle = "none";
                    

                    // Check pathupload
                    if(me.pathupload==""){
                        alert("Set path upload with menu contextuel")
                        return
                    }

                    const files = event.dataTransfer.items;
                    const fileArray = [];

                    // Get all fileEntry
                    for (let i = 0; i < files.length; i++) {
                        const entry = files[i].webkitGetAsEntry();
                        fileArray.push(entry);
                    }

                    // Get all fileEntry
                    processFiles(fileArray);
                    function processFiles(files) {
                        for (let i = 0; i < files.length; i++) {
                            const entry = files[i];
                            if (entry.isFile) {
                                readFile(entry);
                            } else if (entry.isDirectory) {
                                readDirectory(entry.createReader());
                            }
                        }
                    }

                    // Read dir
                    function readDirectory(reader) {
                        reader.readEntries(function(entries) {
                            for (let i = 0; i < entries.length; i++) {
                                const entry = entries[i];
                                if (entry.isFile) {
                                    readFile(entry);
                                } else if (entry.isDirectory) {
                                    readDirectory(entry.createReader());
                                }
                            }
                        });
                    }


                    // Get dir of file
                    function getDir(filePath) {
                        const parts = filePath.split('/');
                        parts.pop();
                        let folderPath = parts.join('/');               
                        return folderPath;
                    }
                    
                    // send file
                    function readFile(fileEntry) {
                        fileEntry.file(function(file) {
                            const formData = new FormData();
                            formData.append('file', file);

                            console.log(fileEntry);

                            // Exemple avec Fetch :
                            fetch('/uploadFile?sessionNumber='+me.terminalVue.sessionNumber+"&pathDir="+
                                        encodeURI(me.pathupload)+
                                        "/"+getDir(fileEntry.fullPath), {
                                method: 'POST',
                                body: formData
                            })
                            .then(response => {
                                if (response.ok) {
                                    console.log('Fichier envoyé avec succès');
                                    this.terminalVue.$refs.info.html("File send successfull.");
                                } else {
                                    this.terminalVue.$refs.info.html("File not send .");
                                }
                            });
                            
                        });
                    }
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
        return         [
            {
                name: 'Download path',
                slug: 'add-star',
            },
            {
                name: 'Set Dir Upload',
                slug: 'remove-star',
            },
            {
                name: 'Copy',
                slug: 'copy',
            }
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
            case "Copy":
                navigator.clipboard.writeText(text);
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
            Id: dataForm.Id,
            Name: dataForm.name,
            URL:`${encodedSchema}://${encodedAuthority}`
        }
    }


    // Get File
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

