import { createApp } from 'vue';
import App from './App.vue';
import router from './router';
import axios from 'axios';
import './index.css'; // Import Tailwind CSS

const app = createApp(App);

app.use(router);

// Set up a base URL for axios
axios.defaults.baseURL = 'http://localhost:8080'; // Ensure this matches your backend URL

app.mount('#app');