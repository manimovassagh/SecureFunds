<template>
  <div class="max-w-md mx-auto mt-10">
    <h2 class="text-2xl font-bold mb-4">Login</h2>
    <form @submit.prevent="login" class="bg-white shadow-md rounded-lg p-6">
      <div class="mb-4">
        <label for="username" class="block text-sm font-medium text-gray-700">Username:</label>
        <input type="text" v-model="username" class="mt-1 block w-full border border-gray-300 rounded-md shadow-sm" required />
      </div>
      <div class="mb-4">
        <label for="password" class="block text-sm font-medium text-gray-700">Password:</label>
        <input type="password" v-model="password" class="mt-1 block w-full border border-gray-300 rounded-md shadow-sm" required />
      </div>
      <button type="submit" class="bg-blue-500 text-white px-4 py-2 rounded">Login</button>
    </form>
    <p class="mt-4 text-center">
      Not registered? <router-link to="/signup" class="text-blue-500">Signup</router-link>
    </p>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue';
import axios from 'axios';

export default defineComponent({
  name: 'Login',
  data() {
    return {
      username: '',
      password: ''
    };
  },
  methods: {
    async login() {
      try {
        const response = await axios.post('/login', {
          username: this.username,
          password: this.password
        });
        console.log('Login response:', response.data);
        const token = response.data.token;
        localStorage.setItem('token', token);

        this.$router.push('/');
      } catch (error) {
        console.error('Login failed', error);
      }
    }
  }
});
</script>