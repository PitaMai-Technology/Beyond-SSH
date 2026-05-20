<script setup lang="ts">
import { ref } from 'vue'
import SshTerminal from './components/SshTerminal.vue'
import SftpExplorer from './components/SftpExplorer.vue'

const drawerLeft = ref(true)
const drawerRight = ref(true)
const tab = ref('ssh')

// Auth state
const isAuthenticated = ref(false)
const loginHost = ref('localhost')
const loginUser = ref('root')
const loginPass = ref('')
const isLoggingIn = ref(false)
const loginError = ref('')

const handleLogin = async () => {
  if (!loginHost.value || !loginUser.value || !loginPass.value) {
    loginError.value = 'Please fill all fields'
    return
  }
  
  isLoggingIn.value = true
  loginError.value = ''
  
  try {
    const result = (window as any).connectSSH(loginHost.value, loginUser.value, loginPass.value)
    if (result === 'Connected') {
      isAuthenticated.value = true
    } else {
      loginError.value = result
    }
  } catch (err: any) {
    loginError.value = err.message || 'Connection failed'
  } finally {
    isLoggingIn.value = false
  }
}

const handleLogout = () => {
  isAuthenticated.value = false
  loginPass.value = ''
}
</script>

<template>
  <v-app class="bg-black">
    
    <!-- Login Screen -->
    <v-main v-if="!isAuthenticated" class="d-flex align-center justify-center h-100">
      <v-card width="400" class="pa-8 rounded-xl elevation-10" color="transparent" style="border: 1px solid #333;">
        <div class="text-center mb-8">
          <v-icon icon="mdi-laptop" size="48" color="blue-lighten-2" class="mb-4"></v-icon>
          <h2 class="text-h4 font-weight-bold">Beyond-SSH</h2>
          <p class="text-grey mt-2">Connect to your server</p>
        </div>
        
        <v-alert v-if="loginError" type="error" variant="tonal" class="mb-6 rounded-lg" density="compact">{{ loginError }}</v-alert>

        <v-text-field v-model="loginHost" label="Host / IP Address" variant="outlined" density="comfortable" rounded="lg" class="mb-2" prepend-inner-icon="mdi-server" hide-details="auto"></v-text-field>
        <v-text-field v-model="loginUser" label="Username" variant="outlined" density="comfortable" rounded="lg" class="mb-2 mt-4" prepend-inner-icon="mdi-account" hide-details="auto"></v-text-field>
        <v-text-field v-model="loginPass" label="Password" type="password" variant="outlined" density="comfortable" rounded="lg" class="mb-6 mt-4" prepend-inner-icon="mdi-lock" @keyup.enter="handleLogin" hide-details="auto"></v-text-field>
        
        <v-btn block color="blue-lighten-2" size="x-large" rounded="pill" class="font-weight-bold text-white text-none elevation-0" @click="handleLogin" :loading="isLoggingIn">
          Connect
        </v-btn>
      </v-card>
    </v-main>

    <!-- Main App Layout -->
    <template v-else>
      <!-- Left Sidebar (Navigation) -->
      <v-navigation-drawer v-model="drawerLeft" app class="pt-4 px-2" width="280" color="black" style="border-right: 1px solid #2f3336;">
        <div class="mb-6 px-4 d-flex align-center">
          <v-btn icon="mdi-console" variant="text" size="large" color="blue-lighten-2" class="mr-2"></v-btn>
        </div>
        
        <v-list nav class="px-2">
          <v-list-item
            prepend-icon="mdi-home-outline"
            title="Home"
            value="home"
            rounded="pill"
            class="mb-2 nav-item"
            active-class="font-weight-bold"
            :ripple="false"
          ></v-list-item>
          <v-list-item
            prepend-icon="mdi-console"
            title="SSH Terminal"
            value="ssh"
            rounded="pill"
            class="mb-2 nav-item"
            @click="tab = 'ssh'"
            :active="tab === 'ssh'"
            active-color="blue-lighten-2"
            active-class="font-weight-bold bg-grey-darken-4"
            :ripple="false"
          ></v-list-item>
          <v-list-item
            prepend-icon="mdi-folder-network-outline"
            title="SFTP Explorer"
            value="sftp"
            rounded="pill"
            class="mb-2 nav-item"
            @click="tab = 'sftp'"
            :active="tab === 'sftp'"
            active-color="blue-lighten-2"
            active-class="font-weight-bold bg-grey-darken-4"
            :ripple="false"
          ></v-list-item>
        </v-list>
        
        <div class="px-4 mt-4">
          <v-btn block color="blue-lighten-2" size="x-large" rounded="pill" class="font-weight-bold text-white text-none elevation-0">
            Post
          </v-btn>
        </div>
        
        <v-spacer></v-spacer>
        
        <!-- User profile mock -->
        <template v-slot:append>
          <div class="pa-4">
            <v-hover v-slot="{ isHovering, props }">
              <v-list-item
                v-bind="props"
                prepend-avatar="https://randomuser.me/api/portraits/lego/1.jpg"
                :title="loginUser"
                :subtitle="'@' + loginHost"
                rounded="pill"
                class="user-profile-btn py-2"
                :class="isHovering ? 'bg-grey-darken-4' : ''"
                @click="handleLogout"
              >
                <template v-slot:append>
                  <v-icon icon="mdi-dots-horizontal" size="small"></v-icon>
                </template>
              </v-list-item>
            </v-hover>
          </div>
        </template>
      </v-navigation-drawer>

      <!-- Main Content Area -->
      <v-main class="bg-black main-border">
        <v-container fluid class="pa-0 h-100 d-flex flex-column">
          <!-- Header -->
          <v-toolbar color="rgba(0,0,0,0.8)" class="px-4 header-blur" density="compact" style="position: sticky; top: 0; z-index: 10; border-bottom: 1px solid #2f3336;">
            <v-toolbar-title class="font-weight-bold text-h6">
              {{ tab === 'ssh' ? 'SSH Terminal' : 'SFTP Explorer' }}
            </v-toolbar-title>
          </v-toolbar>
          
          <div class="pa-0 flex-grow-1 position-relative">
            <div v-show="tab === 'ssh'" class="h-100 pa-4 pb-16">
              <SshTerminal />
            </div>
            <div v-show="tab === 'sftp'" class="h-100 pa-4 pb-16">
              <SftpExplorer />
            </div>
          </div>
        </v-container>
      </v-main>

      <!-- Right Sidebar (Session/Details) -->
      <v-navigation-drawer v-model="drawerRight" location="right" app width="350" color="black" style="border-left: 1px solid #2f3336;">
        <div class="pa-4">
          <v-text-field
            prepend-inner-icon="mdi-magnify"
            placeholder="Search"
            variant="solo-filled"
            density="compact"
            rounded="pill"
            bg-color="grey-darken-4"
            hide-details
            class="mb-6 search-bar"
          ></v-text-field>

          <v-card variant="flat" class="rounded-xl mb-4 bg-grey-darken-4 pb-2">
            <v-card-title class="font-weight-bold text-h6 pt-4 px-4">What's happening</v-card-title>
            <v-list lines="two" bg-color="transparent" class="pt-0">
              <v-hover v-slot="{ isHovering, props }">
                <v-list-item v-bind="props" :class="isHovering ? 'bg-grey-darken-3' : ''" class="px-4 py-2 transition-swing cursor-pointer">
                  <div class="text-caption text-grey">Session Info</div>
                  <div class="font-weight-bold text-body-1 mt-1">Host: {{ loginHost }}</div>
                  <div class="text-caption text-grey mt-1">Connected just now</div>
                </v-list-item>
              </v-hover>
              <v-hover v-slot="{ isHovering, props }">
                <v-list-item v-bind="props" :class="isHovering ? 'bg-grey-darken-3' : ''" class="px-4 py-2 transition-swing cursor-pointer">
                  <div class="text-caption text-grey">System Status</div>
                  <div class="font-weight-bold text-body-1 mt-1">CPU: 5% | RAM: 1.2GB</div>
                  <div class="text-caption text-grey mt-1">Live metrics</div>
                </v-list-item>
              </v-hover>
            </v-list>
            <v-btn variant="text" color="blue-lighten-2" class="text-none px-4" size="small">Show more</v-btn>
          </v-card>
        </div>
      </v-navigation-drawer>
    </template>
  </v-app>
</template>

<style>
/* Twitter-like Micro-interactions & Styling */
.nav-item {
  transition: background-color 0.2s ease;
}
.nav-item:hover:not(.v-list-item--active) {
  background-color: #1a1a1a !important;
}
.v-list-item__content {
  font-size: 1.25rem !important;
}
.user-profile-btn {
  cursor: pointer;
  transition: background-color 0.2s ease;
}
.header-blur {
  backdrop-filter: blur(12px);
}
.search-bar .v-field--variant-solo-filled {
  box-shadow: none !important;
}
.search-bar .v-field--focused {
  background-color: black !important;
  border: 1px solid #4fc3f7 !important;
}
.main-border {
  border-right: 1px solid #2f3336;
}
.transition-swing {
  transition: 0.3s cubic-bezier(0.25, 0.8, 0.5, 1) !important;
}
</style>
