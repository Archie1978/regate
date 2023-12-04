<template>
    <div class="modal-overlay">
      <div class="modal">
        <img class="check" src="~/assets/check-icon.png" alt="" />
  
      <div class="form">
        <div class="form-row">
          <div>Name:</div>
          <input v-model="name" placeholder="username" />
        </div>
         <div class="form-row">
            <div>Password: </div>
            <input v-model="name" placeholder="password" />
          </div>
  
  
          <div>
            <button @click="valid">Validate</button>
          </div>
  
        </div>
      </div>
      <div class="close">
        <img class="close-img" src="~/assets/close-icon.svg" alt="" />
      </div>
    </div>
  </template>
  


    
<script>
export default {
    name: "ModalLogin",
    components: { },
    props: {},
    methods: {

      valid() {
        //this.$emit('validate',this.data);
        if(this.name=="") {
              return ("Name of connection is empty")
        }

        //this.$ws.saveConnection(dataWS);
        this.$emit('close');
      },
    },
    computed: {
      theme() {
        //return this.$store.state.userSettings.darkMode ? "dark" : "";
        return "dark";
      },
    },
    data: () => ({

    }),
  
  };
  </script>