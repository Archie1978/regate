<template>
    <VNavigationDrawer
        expand-on-hover
        rail
        :width="400"
    >

    <v-list>
        <v-list-item
                :prepend-avatar="avatarImg"
                :title="nameDisplay"
                :subtitle="email"
        />
    </v-list>
    <v-diviser></v-diviser>

    <v-list density="compact" nav >
            <v-list-item prepend-icon="mdi-home-city" title="Home" value="home" @click="clickItem('welcome')"></v-list-item>

            <input :value="text" 
                @input="event => text = event.target.value"
            >
            <v-list  nav :max-height="maxHeight" style="overflow-y: scroll;">

              <!--
              <v-list-item  v-for="server in servers"  ref="listServer"
                :key="server.id"
                prepend-icon="mdi-server" 
                :title="server.name"
                :value="server.id"
                @click="clickItemServer(server)"
              >
              </v-list-item>
            -->
            <v-list-item dense v-for="server in servers"  ref="listServer" prepend-icon="mdi-server" :value="server.id" v-bind:key="server.id" >
              <v-row align="center" justify="space-between" >
                <!-- List of servers reported by the app via the websocket -->
                <v-list-item-content>
                  <v-list-item-title class="item-margin" @click="clickItemServer(server)" >{{ server.name }}
                  </v-list-item-title>
                </v-list-item-content>
                <v-list-item-action  class="ml-auto">
                  <!-- Add the button inside the list item -->
                  <v-btn icon @click="modifyItemServer(server)">
                    <v-icon color="red darken-2">mdi-pencil</v-icon>
                  </v-btn>
                  <v-btn icon @click="deleteItemServer(server)">
                    <v-icon color="red darken-2">mdi-delete</v-icon>
                  </v-btn>
                </v-list-item-action>
              </v-row>
            </v-list-item>    
          </v-list>

            <!--
            <v-list  nav :max-height="maxHeight" style="overflow-y: scroll;">
                <v-container fluid fill-height>
                    <v-list-item prepend-icon="mdi-server" title="Servers1" value="servers"></v-list-item>
                    <v-list-item prepend-icon="mdi-server" title="Servers2" value="servers"></v-list-item>
                    <v-list-item prepend-icon="mdi-server" title="Servers3" value="servers"></v-list-item>
                    <v-list-item prepend-icon="mdi-server" title="Servers3" value="servers"></v-list-item>
                    <v-list-item prepend-icon="mdi-server" title="Servers3" value="servers"></v-list-item>
                    <v-list-item prepend-icon="mdi-server" title="Servers3" value="servers"></v-list-item>
                    <v-list-item prepend-icon="mdi-server" title="Servers3" value="servers"></v-list-item>
                    <v-list-item prepend-icon="mdi-server" title="Servers3" value="servers"></v-list-item>
                    <v-list-item prepend-icon="mdi-server" title="Servers3" value="servers"></v-list-item>
                    <v-list-item prepend-icon="mdi-server" title="Servers3" value="servers"></v-list-item>
        
                    <v-list-item prepend-icon="mdi-server" title="Servers3" value="servers"></v-list-item>
                    <v-list-item prepend-icon="mdi-server" title="Servers3" value="servers"></v-list-item>
                    <v-list-item prepend-icon="mdi-server" title="Servers3" value="servers"></v-list-item>
                </v-container>
            </v-list>

            -->
            <v-list-item prepend-icon="mdi-server-plus" title="Add Connection" value="add_conection" @click="clickItem('addConnexion')"></v-list-item>
            <v-list-item prepend-icon="mdi-settings-helper" title="Settings" value="settings" @click="clickItem('settings')"></v-list-item>
            <!--
            <v-list-item prepend-icon="mdi-connection" title="Directly" value="directly" @click="clickItem('directly')"></v-list-item>
            <v-list-item prepend-icon="mdi-exit-to-app" title="Exit" value="exit" ></v-list-item>
            -->

    </v-list>
    </VNavigationDrawer>
</template>

<style>
.item-margin{
  margin-left: 18px;
}

.text-align {
  text-align: "left"
}
.slidebarButton {
  font-size: 12px; /* Taille de police */
  padding: 6px 6px; /* Rembourrage du bouton */
}
</style>


<script>
import eventBus from '../eventBus'


import loadAvatar from '@/assets/avatar.jpg'

export default {
  name: 'SlideBar',

  data: () => ({
    "servers":[],
    avatarImg:loadAvatar,
    nameDisplay:"-----",
    email:""
  }),
  components: {},
  computed:{
    maxHeight:function(){
        var size=window.innerHeight-260;
        return  size+"px";
    }
  },
  methods: {
    clickItem(type){
      
      switch(type){

        case "addConnexion":
          eventBus.$emit('ShowModalConfigurationConnection',{});
          break;

        case "directly":
          eventBus.$emit('ShowModalConfigurationConnection',{});
          break;

        case "welcome":
            this.$ws.tabPanel.addTabCustom("Welcome","Welcome",null);
            break;
        case "settings":
            this.$ws.tabPanel.addTabCustom("Settings","Settings",null);
            break;

        default:
          console.log("item:",type,"not found")
          break;
      }
    },
    
    // Click from item server
    clickItemServer:function(menuItem){
        console.log("menu click:",this,menuItem,menuItem.serverID)
        this.$ws.createSessionMenu=menuItem;
        this.$ws.getStartSession(menuItem.typeProcotol,"Id:"+menuItem.serverID)
    },

    // Modify connexion server
    modifyItemServer:function(server){
      console.log("modify:",server);
      eventBus.$emit('ShowModalConfigurationConnection',server);
    },

    // Delete connexion
    deleteItemServer:function(server){
      console.log("Delete",server)
      this.$ws.deleteConnection(server.id)
    },

    // Create menu into list
    updateMenu:function(rootDataMenu,arraydataWS){
      if(!Array.isArray(arraydataWS)){
        console.log("Nothing group")
        return
      }
      var arraydataWSLength = arraydataWS.length;
      for (var i = 0; i < arraydataWSLength; i++) {
        // Display group
        var data=arraydataWS[i];
        if(data.ServerGroupChildren){
          if(data.ServerGroupChildren.length){
            this.updateMenu(rootDataMenu,data.ServerGroupChildren);
          }
        }
        // display serveur
        if(data.Servers){
          if(data.Servers.length){
            var arrayserverLength = data.Servers.length;
            for (var j = 0; j < arrayserverLength; j++) {
              var dataServer=data.Servers[j];

              // Get protocol
              var protocol = dataServer.URL.split(":")[0];
              
              console.log("dataServer",dataServer);
              // eslint-disable-next-line
              if(listPlugin[protocol]){
                var objectMenuServer={
                  id:dataServer.ID,
                  URL: dataServer.URL,
                  name:dataServer.Name,
                  typeProcotol:protocol.charAt(0).toUpperCase() + protocol.slice(1),
                  serverID:dataServer.ID,
                  // eslint-disable-next-line
                  protocolClass:listPlugin[protocol],
                }
                rootDataMenu.push(objectMenuServer)
              }
            }
          }
        }
      }
    },

    // Create datastructure for menu tree
    updateMenuTree:function(rootMenu,arraydataWS){
      var arrayChild=new Array();
      var arraydataWSLength = arraydataWS.length;
      for (var i = 0; i < arraydataWSLength; i++) {
        var data=arraydataWS[i];
        var objectMenu={
          href:rootMenu.href+'/'+data.Name,
          name: data.Name,
          children: []
        }
        if(data.ServerGroupChildren){
          if(data.ServerGroupChildren.length){
            console.log("children found")
            this.updateMenuTree(objectMenu,data.ServerGroupChildren)
          }
        }
        // Add serveur into slide bar
        if(data.Servers){
          if(data.Servers.length){
            var arrayserverLength = data.Servers.length;
            for (var j = 0; j < arrayserverLength; j++) {
              var dataServer=data.Servers[j];
              console.log(dataServer);
              var u=new URL("/",dataServer.URL);
              var objectMenuServer={
                href:data.href+'/'+dataServer.Name,
                name: dataServer.Name,
                serverID: dataServer.ID,
                loginType: dataServer.LoginType,
              }

              var protocol=(u.protocol).slice(0, -1);
              // eslint-disable-next-line
              if(listPlugin[protocol]){
                // eslint-disable-next-line
                var objectProtocol=new listPlugin[protocol];
                // eslint-disable-next-line
                objectMenuServer["protocolClass"]=listPlugin[protocol];
                objectMenuServer["typeUI"]=objectProtocol.typeUI;
                objectMenuServer["icon"]=objectProtocol.icon;
                objectMenuServer["typeProcotol"]=protocol.charAt(0).toUpperCase() + protocol.slice(1);
                objectMenu.children.push(objectMenuServer);
              }
            }
          }
        }
        arrayChild.push(objectMenu)
      }
    }
  },
  mounted() {
    var me=this;

    // Refresh List of server
    eventBus.$on("LISTSERVER",function(data){
      var serverList=[];
      me.updateMenu(serverList,data.Msg.ServerGroupChildren);

      // sort list
      serverList.sort((a, b) => {
        const nameA = a.name.toUpperCase(); // ignore upper and lowercase
        const nameB = b.name.toUpperCase(); // ignore upper and lowercase
        if (nameA < nameB) {
          return -1;
        }
        if (nameA > nameB) {
          return 1;
        }

        // names must be equal
        return 0;
      });
      me.$data["servers"]=serverList;
      me.$forceUpdate();
    })

    // Refresh List of server
    eventBus.$on("USER_PROFILE",function(object){
      console.log("ON bus",object);
      console.log("loadAvatar",loadAvatar)
      me.avatarImg= "data:image/png;base64,"+object.Avatar;
      me.name=object.Name;
      me.lastName=object.LastName;
      me.nameDisplay=object.Name+" "+object.LastName;
      me.email=object.Email;                    
    })

    // Connect to server regate and get List of server
    this.$ws.connected();
    this.$ws.getListMenu();
  }
}
</script>
