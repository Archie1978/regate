
import eventBus from './eventBus'


export default class WebSocketRemoteDesktop {
 /*
    constructor() {
        this.connection = null
        this.eventBus = new eventBus()
    }
*/
    connected(){
        console.log("Starting connection to WebSocket Server")
        console.log(process.env.DEVELOPPEMENT)

        var me=this;
        var url="ws://"+location.hostname+":8088/ws"
        if ( process.env.URLWS ) {
            url=process.env.URLWS
        }
      
        this.connection = new WebSocket(url);
        this.connection.onmessage = function(event) {
            var object=JSON.parse(event.data)
            console.info("recept:",object);
            switch(object.Command.toUpperCase()){
                case "STARTED":
                    //Create tab
                    console.log(me,me.createSessionMenu.protocolClass);
                    me.tabPanel.addTabRemote(
                        object.Session,
                        me.createSessionMenu.name,
                        "index"+object.Session,
                        new me.createSessionMenu.protocolClass,
                    );

                    // Get End session
                    eventBus.$on(object.Session+"_End", () => {
                        console.log("End session")
                        me.tabPanel.removeTab(
                            "index"+object.Session,
                        );
                    });

                    break;
                case "ERROR":
                    alert(object.Msg);
                    break;
                    
                default:
                    if(object.Session==0){
                        eventBus.$emit(object.Command, object);
                    }else{
                        eventBus.$emit(object.Session+"_"+object.Command, object);
                    }
                    break;
            }

        }

        this.connection.onopen = function(event) {
            console.log("Successfully connected to the echo websocket server...")
            console.log(event);
        }

        this.connection.onerror = function(event) {
            console.log("error web socket",event)
        }

        this.connection.onclose =function () {
            console.log('La connexion a été fermée avec succès.');
            eventBus.$emit("WebSocketLost", {});
        }
    
    }

    login(user,password){
        console.log("connect",user,password)
        //connected()  
    }

    getListMenu(){
        this.sendMessage({Command:"LISTSERVER" })
    }
    getVersion(){
        this.sendMessage({Command:"VERSION" })
    }
    saveConnection(recordServer){
        this.sendMessage({Command:"SaveConnection", Msg:recordServer})
    }
    deleteConnection(recordServerID){
        console.log("lll");
        this.sendMessage({Command:"DeleteConnection", Msg:recordServerID})
        console.log("lll2");
    }

    // Start session with option
    // Id:int => serverID
    // object => connection with protocol
    getStartSession(typeProtocol,option){
        this.sendMessage({Command:"START",TypeProtocol:typeProtocol, Msg: option})
    }

    sendMessage(object){
        var me=this;
        console.log(object,this);
        if(!this.connection.readyState){
            console.log("Connexion not ready");
            setTimeout(function () {
                me.sendMessage(object)
            }, 1.0 * 1000);
        }else{
            console.log("Send:",this,object);
            this.connection.send(JSON.stringify(object));
        }
    }
}



/*


import EventEmitter from 'events'


export default class WebSocketRemoteDesktop extends EventEmitter{
 
    constructor() {
        super()
        this.connection = null
    }
    created() {
        
    }

    connected(){
        console.log("Starting connection to WebSocket Server")
        console.log(process.env.DEVELOPPEMENT)
        var me=this;
        var url="ws://"+location.hostname+":8088/ws"
        if ( process.env.URLWS ){
            url=process.env.URLWS
        }
      
        this.connection = new WebSocket(url);
        this.connection.onmessage = function(event) {
            var object=JSON.parse(event.data)
            me.emit(object.Command, object)
            console.log(event);

        }

        this.connection.onopen = function(event) {
            console.log("Successfully connected to the echo websocket server...")
            console.log(event);
        }

        this.connection.onerror = function(event) {
            console.log("erro web socket")
            console.log(event);
        }
    
    }

    login(user,password){
        console.log("connect",user,password)
        //connected()  
    }

    getListMenu(){
        this.sendMessage({Command:"LISTSERVER" })
    }
    getVersion(){
        this.sendMessage({Command:"VERSION" })
    }

    sendMessage(object){
        var me=this;
        if(!this.connection.readyState){
            console.log("Connexion not ready");
            setTimeout(function () {
                me.sendMessage(object)
            }, 1.0 * 1000);
        }else{
            console.log("Send:",object);
            this.connection.send(JSON.stringify(object));
        }
    }

}

*/
