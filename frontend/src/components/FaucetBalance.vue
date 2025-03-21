<template>
  <div class="bg-[#141414] rounded-lg p-6 flex items-center justify-between">
    <div class="flex items-center">
      <div class="bg-[#00ffa3] bg-opacity-10 rounded-full p-3">
        <svg class="h-6 w-6 text-[#00ffa3]" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
      </div>
      <div class="ml-4">
        <h2 class="text-lg font-medium text-[#405045]">Faucet Balance</h2>
        <div class="flex items-center mt-1">
          <p class="text-2xl font-bold text-[#00ffa3]">{{ formatBalance(balance) }} SOL</p>
          <span v-if="isLoading" class="ml-2">
            <svg class="animate-spin h-5 w-5 text-[#00ffa3]" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
          </span>
        </div>
      </div>
    </div>
    <button 
      @click="fetchBalance" 
      class="text-[#405045] hover:text-[#00ffa3] transition-colors duration-200"
      :disabled="isLoading"
    >
      <svg class="h-5 w-5" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor">
        <path fill-rule="evenodd" d="M4 2a1 1 0 011 1v2.101a7.002 7.002 0 0111.601 2.566 1 1 0 11-1.885.666A5.002 5.002 0 005.999 7H9a1 1 0 010 2H4a1 1 0 01-1-1V3a1 1 0 011-1zm.008 9.057a1 1 0 011.276.61A5.002 5.002 0 0014.001 13H11a1 1 0 110-2h5a1 1 0 011 1v5a1 1 0 11-2 0v-2.101a7.002 7.002 0 01-11.601-2.566 1 1 0 01.61-1.276z" clip-rule="evenodd" />
      </svg>
    </button>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import axios from 'axios'
import { apiBaseUrl } from '../config'

const balance = ref(0)
const isLoading = ref(false)

const formatBalance = (value) => {
  if (!value && value !== 0) return '0.00'
  return Number(value).toFixed(2)
}

const fetchBalance = async () => {
  if (isLoading.value) return
  
  isLoading.value = true
  try {
    const response = await axios.get(`${apiBaseUrl}/api/balance`)
    balance.value = response.data.balance || 0
  } catch (error) {
    console.error('Error fetching balance:', error)
  } finally {
    isLoading.value = false
  }
}

onMounted(() => {
  fetchBalance()
})
</script> 