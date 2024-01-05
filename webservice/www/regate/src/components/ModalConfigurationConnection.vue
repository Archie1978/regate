<template>
  <div class="modal-overlay">
    <div class="modal">
      <img class="check" src="~/assets/check-icon.png" alt="" />

    <div class="form">
      <div class="form-row">
        <div>Connection Name:</div>
        <input v-model="name" placeholder="edit me" />
      </div>
        <div class="form-row">
          <div>Protocol: </div>
          <select v-model="protocol">
            <option disabled value="">Please select one</option>
            <option value="ssh">Ssh ( Posix/Linux/Unix )</option>
            <option value="rdp">Rdp ( Windows )</option>
          </select>
        </div>

        <div v-show="protocol=='ssh'">
          <div class="form-row">
            <div>Host:</div>
            <input v-model="host" placeholder="sd20932.online.net/192.168.10.22" />
          </div>

          <div class="form-row">
            <div>Port:</div>
            <input v-model="port" placeholder="22" />
          </div>
          <div class="form-row">
            <div>User:</div>
            <input v-model="user" placeholder="foouser" />
          </div>

          <div class="form-row">
            <div>Password:</div>
            <input v-model="password" type="password" placeholder="My password" />
          </div>
        </div>

        <div v-show="protocol=='rdp'">
          <div class="form-row">
            <div>Host:</div>
            <input v-model="host" placeholder="sd20932.online.net/192.168.10.22" />
          </div>
          <div class="form-row">
            <div>User:</div>
            <input v-model="user" placeholder="foouser" />
          </div>

          <div class="form-row">
            <div>Domain:</div>
            <input v-model="domain" placeholder="microsoft" />
          </div>

          <div class="form-row">
            <div>Password:</div>
            <input v-model="password" type="password" placeholder="edit me" />
          </div>
        </div>

        <div>
          <button @click="closeDialog">Cancel</button>
          <button @click="valid">Validate</button>
        </div>

      </div>
    </div>
    <div class="close">
      <img class="close-img" src="~/assets/close-icon.svg" alt="" />
    </div>
  </div>
</template>

<style scoped>

.modal-overlay {
  position: fixed;
  top: 0;
  bottom: 0;
  left: 0;
  right: 0;
  display: flex;
  justify-content: center;
  background-color: #000000da;
}

.modal {
  text-align: center;
  background-color: white;
  height: 500px;
  width: 500px;
  margin-top: 10%;
  padding: 60px 0;
  border-radius: 20px;
}
.close {
  margin: 10% 0 0 16px;
  cursor: pointer;
}

.close-img {
  width: 25px;
}

.check {
  width: 150px;
}

h6 {
  font-weight: 500;
  font-size: 28px;
  margin: 20px 0;
}

p {
  font-size: 16px;
  margin: 20px 0;
}

button {
  background-color: #ac003e;
  width: 150px;
  height: 40px;
  color: white;
  font-size: 14px;
  border-radius: 16px;
  margin-top: 50px;
  margin: 50px 20px 15px 20px;
}


.form .form-row  {
  display: flex;
  flex-direction: row;
  margin: 0px 30px;
}
.form .form-row  :first-child{
  flex: 3;
  padding-left: 10px;
  border-bottom:#ac003e;
  border-bottom-style: solid;
}

.form .form-row  :nth-child(2){
  border-style: double;
}
.form .form-row  * {
  flex: 5;
  margin: 10px 10px;
  text-align: left;
}
</style>

    
<script>

export default {
  name: "ModalConfigurationConnection",
  components: { },
  
  methods: {

    handleShowEvent(){
      console.log("====Show");
    },

    // Close dial
    closeDialog() {
      this.$emit('close');
    },


    // Set form by server recorder
    setRecord(recordServer){
      this.name=recordServer.name;

      // Get protocole
      const urlProtocol = new URL(recordServer.URL);
      this.protocol=urlProtocol.protocol.slice(0, -1);

      // Get other information replace juste first element by http
      let urlStringModified = recordServer.URL.replace(this.protocol, "http");
      let ur1Modified=new URL(urlStringModified);
      this.Id=recordServer.id;

      this.host=decodeURI(ur1Modified.hostname);
      this.user=decodeURI(ur1Modified.username);
      this.password=decodeURI(ur1Modified.password) ;
    },



    // Validate information
    valid() {
      
      //this.$emit('validate',this.data);
      if(this.name=="") {
        alert("name of connection is empty")
        return
      }
      // eslint-disable-next-line
      if(!listPlugin[this.protocol]){
        alert("protocol unknown") ;
        return;
      }

      // Encode configuration into URL by plugin
      // eslint-disable-next-line
      var dataWS=listPlugin[this.protocol].encodeConfiguration(this);
      console.log(typeof dataws)
      if (typeof(dataWS)=== "string") {
        alert(dataWS);
        return;
      }

      this.$ws.saveConnection(dataWS);
      this.$emit('close');
    },
  },

  // Change theme
  computed: {
    theme() {
      //return this.$store.state.userSettings.darkMode ? "dark" : "";
      return "dark";
    },
  },

  // Data form
  data: () => ({
    "Id":0,
    "protocol":"",
    "port":0,
    "name":"",
    "user":"",
    "domain":"",
    "password":"",
  }),

};
</script>