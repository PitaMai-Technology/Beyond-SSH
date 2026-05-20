<script setup lang="ts">
import { ref } from 'vue'

interface FileNode {
  name: string
  isDir: boolean
  size?: string
  children?: FileNode[]
}

const fileTree = ref<FileNode[]>([])
const isLoading = ref(false)
const error = ref('')

const loadDir = async (path: string = '.') => {
  isLoading.value = true
  error.value = ''
  try {
    const resultStr = (window as any).listSFTPDirectory(path)
    if (resultStr && resultStr.startsWith('{"error"')) {
      const errObj = JSON.parse(resultStr)
      error.value = errObj.error
    } else if (resultStr) {
      fileTree.value = JSON.parse(resultStr)
    }
  } catch (e: any) {
    error.value = e.message || 'Error calling listSFTPDirectory'
  } finally {
    isLoading.value = false
  }
}
</script>

<template>
  <v-card variant="outlined" class="h-100 d-flex flex-column rounded-lg">
    <v-toolbar color="transparent" density="compact" class="border-b px-2">
      <v-btn icon="mdi-refresh" variant="text" size="small" @click="loadDir('.')"></v-btn>
      <v-btn icon="mdi-folder-plus" variant="text" size="small" :disabled="!fileTree.length"></v-btn>
      <v-btn icon="mdi-file-plus" variant="text" size="small" :disabled="!fileTree.length"></v-btn>
      <v-spacer></v-spacer>
      <v-btn icon="mdi-cloud-upload" variant="text" size="small" color="primary" :disabled="!fileTree.length"></v-btn>
    </v-toolbar>
    
    <div class="flex-grow-1 overflow-y-auto pa-2">
      <div v-if="isLoading" class="text-center pa-4 text-grey">Loading...</div>
      <div v-else-if="error" class="text-error pa-4">{{ error }}</div>
      <div v-else-if="fileTree.length === 0" class="text-center pa-4 text-grey">
        Click refresh to load directory. (Must be connected to SSH first)
      </div>
      <v-list density="compact" nav v-else>
        <template v-for="node in fileTree" :key="node.name">
          <v-list-group v-if="node.isDir" :value="node.name">
            <template v-slot:activator="{ props }">
              <v-list-item
                v-bind="props"
                prepend-icon="mdi-folder"
                :title="node.name"
                color="primary"
              ></v-list-item>
            </template>
            
            <template v-for="child in node.children" :key="child.name">
               <v-list-item
                 v-if="!child.isDir"
                 prepend-icon="mdi-file-document-outline"
                 :title="child.name"
                 :subtitle="child.size"
                 class="pl-8"
               ></v-list-item>
               <v-list-group v-else :value="child.name">
                 <template v-slot:activator="{ props }">
                   <v-list-item
                     v-bind="props"
                     prepend-icon="mdi-folder"
                     :title="child.name"
                     class="pl-8"
                   ></v-list-item>
                 </template>
                 <template v-for="grandchild in child.children" :key="grandchild.name">
                   <v-list-item
                     prepend-icon="mdi-file-document-outline"
                     :title="grandchild.name"
                     :subtitle="grandchild.size"
                     class="pl-12"
                   ></v-list-item>
                 </template>
               </v-list-group>
            </template>
          </v-list-group>
          
          <v-list-item
            v-else
            prepend-icon="mdi-file-document-outline"
            :title="node.name"
            :subtitle="node.size"
          ></v-list-item>
        </template>
      </v-list>
    </div>
  </v-card>
</template>

<style scoped>
.border-b {
  border-bottom: 1px solid rgba(var(--v-border-color), var(--v-border-opacity));
}
</style>
