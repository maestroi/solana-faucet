<template>
  <div class="h-screen w-full flex flex-col bg-[#141414] overflow-hidden">
    <!-- Header -->
    <header class="w-full bg-[#191b1c] border-b border-[#405045] flex-shrink-0">
      <div class="w-full px-4 sm:px-6 lg:px-8">
        <div class="flex items-center justify-between h-16">
          <div class="flex items-center">
            <img src="/solana-logo.svg" alt="Solana Logo" class="h-8 w-8 mr-3" />
            <h1 class="text-xl font-bold text-[#00ffa3]">Solana Faucet</h1>
          </div>
          <div class="flex items-center">
            <div class="flex items-center">
              <div class="relative">
                <div class="h-3 w-3 rounded-full" :class="isHealthy ? 'bg-green-500' : 'bg-red-500'"></div>
                <div class="absolute h-3 w-3 rounded-full animate-ping" :class="isHealthy ? 'bg-green-500' : 'bg-red-500'"></div>
              </div>
              <span class="ml-2 text-sm text-[#405045]">{{ isHealthy ? 'Online' : 'Offline' }}</span>
            </div>
          </div>
        </div>
      </div>
    </header>

    <!-- Main Content -->
    <main class="flex-1 w-full overflow-auto">
      <div class="w-full px-4 sm:px-6 lg:px-8 py-8">
        <slot></slot>
      </div>
    </main>

    <!-- Footer -->
    <footer class="w-full bg-[#191b1c] border-t border-[#405045] flex-shrink-0">
      <div class="w-full px-4 sm:px-6 lg:px-8">
        <div class="flex items-center justify-between h-16">
          <div class="text-sm text-[#405045]">
            © {{ new Date().getFullYear() }} Solana Faucet
          </div>
          <div class="flex items-center space-x-6">
            <a href="https://github.com/maestroi/solana-faucet/blob/main/readme.md" target="_blank" class="text-sm text-[#405045] hover:text-[#00ffa3] transition-colors duration-200">
              Documentation
            </a>
            <a href="https://github.com/maestroi" target="_blank" class="text-sm text-[#405045] hover:text-[#00ffa3] transition-colors duration-200">
              GitHub
            </a>
          </div>
        </div>
      </div>
    </footer>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import axios from 'axios'
import { apiBaseUrl } from '../config'

const isHealthy = ref(false)

const checkHealth = async () => {
  try {
    const response = await axios.get(`${apiBaseUrl}/api/health`)
    isHealthy.value = response.data.ok === true
  } catch (error) {
    console.error('Health check failed:', error)
    isHealthy.value = false
  }
}

onMounted(() => {
  checkHealth()
  // Check health every 30 seconds
  setInterval(checkHealth, 30000)
})
</script> 