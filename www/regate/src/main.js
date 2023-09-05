import { createApp } from 'vue'
import App from './App.vue'
import vuetify from './plugins/vuetify'
import { loadFonts } from './plugins/webfontloader'

loadFonts()

const app=createApp(App);

// Globals variables
import WebSocketRemoteDesktop from "./WebSocketRemoteDesktop.js";
app.config.globalProperties.$ws =new WebSocketRemoteDesktop();

app.use(vuetify)
    .mount('#app');