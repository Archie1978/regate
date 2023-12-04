<template>
  <div class="modal-overlay">
    <div class="modal">
      <img class="check" src="~/assets/check-icon.png" alt="" />
      <h6>{{ title }}</h6>
      <p>{{ message }}</p>
      <button @click="close">OK</button>
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
  height: 300px;
  width: 450px;
  margin-top: 10%;
  padding: 10px 0;
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
}
</style>

    
<script>
  
  export default {
    name: "ModalMessage",
    components: { },
    props: {
      title: String,
      message: String, 
    },
    methods: {
      close() {
        this.$emit('close');
      },
    },
    computed: {
      theme() {
        //return this.$store.state.userSettings.darkMode ? "dark" : "";
        return "dark";
      },
    }
  };
  </script>