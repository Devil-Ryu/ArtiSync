import {createApp} from 'vue'
import App from './App.vue'
import TDesign from 'tdesign-vue-next';
import { createPinia } from 'pinia';
import router from './router';
import './style.css';

import  'tdesign-vue-next/es/style/index.css';
const pinia = createPinia()
const app = createApp(App)
app.use(TDesign)
app.use(pinia)
app.use(router)
app.mount('#app')
