import { createApp } from 'vue';
import App from './App.vue';
import vuetify from './plugins/vuetify'; // Import Vuetify plugin

createApp(App)
  .use(vuetify) // Use Vuetify
  .mount('#app');
