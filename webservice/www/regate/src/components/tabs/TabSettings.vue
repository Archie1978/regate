<template>
  <v-row>
    <div class="form">
      <div class="settings-title">Security</div>
      <v-row class="settings-index">
        <v-col>
          <div>Finger Host (known_host)</div>
          <div class="form-row">
              <div>Activate:</div>
              <input v-model="kh_activate" type="checkbox" placeholder="Activate security server" />
          </div>
          <div class="form-row">
              <div>Check By DNS ( NOT YET) :</div>
              <input v-model="kh_dns" type="checkbox" placeholder="Check" />
          </div>
          <div class="form-row">
              <div>Text:</div>
              <textarea class="kh_list"  v-model="kh_list" type="text" placeholder="" />
          </div>
        </v-col>
      </v-row>
      <hr>
      <v-row class="settings-index">
        <v-col>
          <div>Certificates Roots/Intermedaires</div>
          <div class="form-row">
              <div>Activate:</div>
              <input v-model="cert_activate" type="checkbox" placeholder="Activate security server" />
          </div>
          <div class="form-row">
              <div>Text:</div>
              <textarea class="cert_list" v-model="cert_list" type="text" placeholder="" />
          </div>
        </v-col>
     </v-row> 
     <row>
          <button @click="valid">Validate</button>
    </row>
    </div>
  </v-row>
</template>


<style scoped>

.settings-title{
  font-size: 36px;
  font-weight: bold;
  margin-bottom: 10px;
}

.settings-index {
  margin-left: 10px;
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

.form {
  margin-left:20px;
}

.form > .v-row {
  margin-bottom: 50px;
}

.form .form-row  {
  display: flex;
  flex-direction: row;
  margin: 0px 30px;
  width: 1100px;
}
.form .form-row  :first-child{
  flex: 3;
  padding-left: 10px;
  border-bottom:#ac003e;
  border-bottom-style: solid;
  align-self: end;
}

.form .form-row  :nth-child(2){
  border-style: double;
}
.form .form-row  * {
  flex: 5;
  margin: 10px 10px;
  text-align: left;
}

.kh_list {
  height: 183px;
  font-size: small;
}

.cert_list {
  height: 303px;
  font-size: small;
}
</style>

<script>
import eventBus from '../../eventBus'

export default {
    name: "TabSettings",
    components: {},

  // Set form by server recorder
  mounted() {
    
    // Get data
    var me=this;
    eventBus.$on("SETTING_SECURITY", (messageWS) => {

      // update data from WS
      for (let propname in messageWS) {
        me[propname.toLowerCase()] = messageWS[propname];
      }
    });

    this.$ws.getSettingSecurity()
  },

  data(){
    return {
      kh_activate:false,
      kh_dns:false,
      kh_list:"",
      cert_activate:false,
      cert_list:"",
    }
  },

  methods: {

    // Validate information
    valid() {
      this.$ws.saveSettingSecurity({
          kh_activate: this.kh_activate,
          kh_dns: this.kh_dns,
          kh_list: this.kh_list,
          cert_activate: this.cert_activate,
          cert_list: this.cert_list
      });
      //this.$emit('close');
    },
  },

    props: {
      disabled: {
        Type: Boolean,
        required: true,
        default: false,
      },
      typePanel: {
        Type: String,
        required: true,
        default: "",
      },
    }
}
</script>