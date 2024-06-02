<template>
  <div class="max-w-md mx-auto mt-10">
    <h1 class="text-3xl font-bold mb-6 text-center">Welcome to Secure Funds</h1>
    <div v-if="user">
      <div class="flex justify-between items-center mb-6">
        <h2 class="text-2xl font-bold">Welcome, {{ user.username }}</h2>
        <button @click="logout" class="bg-red-500 text-white px-4 py-2 rounded">Logout</button>
      </div>
      <div v-if="account" class="bg-white shadow-md rounded-lg p-6">
        <h3 class="text-xl font-bold mb-4">Account Balance: {{ account.balance }}</h3>
        <h3 class="text-xl font-bold mb-2">Account History</h3>
        <ul>
          <li v-for="transaction in transactions" :key="transaction.id" class="border-b py-2">
            <span class="font-semibold">{{ transaction.type }}:</span> {{ transaction.amount }}
          </li>
        </ul>
      </div>
    </div>
    <div v-else>
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
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue';
import axios from 'axios';

interface User {
  id: number;
  username: string;
}

interface Transaction {
  id: number;
  type: string;
  amount: number;
}

interface Account {
  id: number;
  balance: number;
  transactions: Transaction[];
}

export default defineComponent({
  name: 'Home',
  data() {
    return {
      user: null as User | null,
      account: null as Account | null,
      transactions: [] as Transaction[],
      username: '',
      password: ''
    };
  },
  async created() {
    await this.checkAuthentication();
  },
  methods: {
    async checkAuthentication() {
      try {
        const token = localStorage.getItem('token');
        if (token) {
          console.log('Token found:', token);

          // Fetch account history directly if account id is known
          const accountId = 1; // Replace with actual account id if available
          const historyResponse = await axios.get(`/account/${accountId}/history`, {
            headers: {
              Authorization: `Bearer ${token}`
            }
          });
          console.log('History response:', historyResponse.data);
          this.transactions = historyResponse.data.transactions;

          // Assuming balance is returned as part of account history response
          this.account = { id: accountId, balance: historyResponse.data.balance, transactions: this.transactions };
          this.user = { id: historyResponse.data.user_id, username: historyResponse.data.username };
        }
      } catch (error) {
        console.error('Failed to fetch account data', error);
        this.user = null;
      }
    },
    async login() {
      try {
        const response = await axios.post('/login', {
          username: this.username,
          password: this.password
        });
        console.log('Login response:', response.data);
        const token = response.data.token;
        localStorage.setItem('token', token);

        this.user = { id: response.data.user_id, username: this.username }; // Assuming login response includes user_id and username

        // Fetch account history directly after login
        const accountId = 1; // Replace with actual account id if available
        const historyResponse = await axios.get(`/account/${accountId}/history`, {
          headers: {
            Authorization: `Bearer ${token}`
          }
        });
        console.log('History response after login:', historyResponse.data);
        this.transactions = historyResponse.data.transactions;

        // Assuming balance is returned as part of account history response
        this.account = { id: accountId, balance: historyResponse.data.balance, transactions: this.transactions };

        this.$router.push('/');
      } catch (error) {
        console.error('Login failed', error);
      }
    },
    logout() {
      localStorage.removeItem('token');
      this.user = null;
      this.account = null;
      this.transactions = [];
      this.$router.push('/login');
    }
  }
});
</script>