<script setup>
import { useRouter } from 'vue-router'
import { useAuthStore } from '../store/auth'

const router = useRouter()
const authStore = useAuthStore()

const handleLogout = () => {
  authStore.logout()
  router.push({ name: 'Login' })
}
</script>

<template>
  <div class="min-h-screen bg-gray-100 flex flex-col md:flex-row">
    <!-- Sidebar -->
    <aside class="w-full md:w-64 bg-gray-900 text-white flex-shrink-0">
      <div class="p-6">
        <h2 class="text-2xl font-bold text-blue-400">Panel Admin</h2>
        <p class="text-sm text-gray-400 mt-1">CBT SMP</p>
      </div>
      
      <nav class="mt-2 space-y-1 px-3">
        <router-link :to="{ name: 'AdminDashboard' }" class="block px-3 py-2 rounded-md text-base font-medium hover:bg-gray-800 hover:text-white" exact-active-class="bg-blue-600 text-white">Dashboard</router-link>
        <router-link :to="{ name: 'AdminGuru' }" class="block px-3 py-2 rounded-md text-base font-medium hover:bg-gray-800 hover:text-white" exact-active-class="bg-blue-600 text-white">Data Guru</router-link>
        <router-link :to="{ name: 'AdminSiswa' }" class="block px-3 py-2 rounded-md text-base font-medium hover:bg-gray-800 hover:text-white" exact-active-class="bg-blue-600 text-white">Data Siswa</router-link>
      </nav>
    </aside>

    <!-- Main Content -->
    <div class="flex-1 flex flex-col overflow-hidden">
      <!-- Top header -->
      <header class="bg-white shadow-sm flex items-center justify-between px-6 py-4">
        <h2 class="text-xl font-semibold text-gray-800">CBT Administrator</h2>
        <div class="flex items-center gap-4">
          <span class="text-sm text-gray-600">Halo, {{ authStore.user?.username }}</span>
          <button @click="handleLogout" class="text-sm bg-red-50 text-red-600 px-3 py-1.5 rounded-md hover:bg-red-100 font-medium transition">
            Logout
          </button>
        </div>
      </header>

      <!-- Main viewport -->
      <main class="flex-1 overflow-x-hidden overflow-y-auto bg-gray-100 p-6">
        <router-view></router-view>
      </main>
    </div>
  </div>
</template>
