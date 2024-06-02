<template>
  <div class="max-w-md mx-auto mt-10">
    <h2 class="text-2xl font-bold mb-4">Signup</h2>
    <form @submit.prevent="signup" class="bg-white shadow-md rounded-lg p-6">
      <div class="mb-4">
        <label for="username" class="block text-sm font-medium text-gray-700">Username:</label>
        <input type="text" v-model="username" class="mt-1 block w-full border border-gray-300 rounded-md shadow-sm" required />
      </div>
      <div class="mb-4">
        <label for="password" class="block text-sm font-medium text-gray-700">Password:</label>
        <input type="password" v-model="password" class="mt-1 block w-full border border-gray-300 rounded-md shadow-sm" required />
      </div>
      <button type="submit" class="bg-blue-500 text-white px-4 py-2 rounded">Signup</button>
    </form>
    <p class="mt-4 text-center">
      Already registered? <router-link to="/login" class="text-blue-500">Login</router-link>
    </p>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue';
import axios from 'axios';

export default defineComponent({
  name: 'Signup',
  data() {
    return {
      username: '',
      password: ''
    };
  },
  methods: {
    async signup() {
      try {
        const response = await axios.post('/signup', {
          username: this.username,
          password: this.password
        });
        console.log('Signup response:', response.data);
        this.$router.push('/login');
      } catch (error) {
        console.error('Signup failed', error);
      }
    }
  }
});
</script>