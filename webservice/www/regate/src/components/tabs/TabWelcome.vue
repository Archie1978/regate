<template>
  <v-row>
    <v-col>
        <h1>Regate Application</h1>
    
        <br>
        <hr>
        <br>
        <img :src="require(`@/assets/logo.png`)" height="100">
        <br>
        <h3>Created by</h3>
        <v-list>
            <v-list-item title="Author: VALMIR Christophe"></v-list-item>
            <v-list-item title="Licence: GPLv3"></v-list-item>
            <v-list-item title="Since: 30/04/2023"></v-list-item>
            <v-list-item title="Date version: 03/09/2023"></v-list-item>
        </v-list>

        <h3>Technology</h3>
        <v-list>
            <v-list-item title="Golang (>1.18)"></v-list-item>
            <v-list-item title="Javascript"></v-list-item>

            <v-list>
              <v-list-item title="VueJs"></v-list-item>
              <v-list-item title="TerminalJs"></v-list-item>
            </v-list>
        </v-list>
    </v-col>
</v-row>
</template>



<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
h3 {
  margin: 40px 0 0;
}
ul {
  list-style-type: none;
  padding: 0;
}
li {
  display: inline-block;
  margin: 0 10px;
}
a {
  color: #42b983;
}
</style>


<script>

export default {
    name: "TabCanvas",
    components: {},

}
</script>