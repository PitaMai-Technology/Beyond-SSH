<script setup lang="ts">
import { ref, nextTick } from 'vue'

const output = ref<{type: 'command' | 'output' | 'error', text: string}[]>([
  { type: 'output', text: 'Welcome to Beyond-SSH Terminal.' },
  { type: 'output', text: 'To connect: ssh <host> <user> <pass>' }
])
const currentCommand = ref('')
const terminalBody = ref<HTMLElement | null>(null)
const isConnected = ref(false)

const executeCommand = async () => {
  if (!currentCommand.value.trim()) return
  
  const cmd = currentCommand.value
  output.value.push({ type: 'command', text: `$ ${cmd}` })
  currentCommand.value = ''
  scrollToBottom()

  if (cmd.startsWith('ssh ')) {
    const parts = cmd.split(' ')
    if (parts.length === 4) {
      const [_, host, user, pass] = parts
      output.value.push({ type: 'output', text: `Connecting to ${host} via proxy...` })
      try {
        const result = (window as any).connectSSH(host, user, pass)
        output.value.push({ type: 'output', text: result })
        if (result === 'Connected') isConnected.value = true
      } catch (err: any) {
        output.value.push({ type: 'error', text: `Error: ${err.message}` })
      }
    } else {
      output.value.push({ type: 'error', text: 'Usage: ssh <host> <user> <pass>' })
    }
    scrollToBottom()
    return
  }

  if (!isConnected.value) {
    output.value.push({ type: 'error', text: 'Not connected. Run: ssh <host> <user> <pass>' })
    scrollToBottom()
    return
  }
  
  try {
    const result = (window as any).executeCommand(cmd)
    output.value.push({ type: 'output', text: result })
  } catch (err: any) {
    output.value.push({ type: 'error', text: `Error: ${err.message}` })
  }
  
  scrollToBottom()
}

const scrollToBottom = async () => {
  await nextTick()
  if (terminalBody.value) {
    terminalBody.value.scrollTop = terminalBody.value.scrollHeight
  }
}
</script>

<template>
  <div class="terminal-container d-flex flex-column h-100 rounded-lg bg-black elevation-2">
    <div 
      class="terminal-body flex-grow-1 overflow-y-auto pa-4 font-mono text-body-2"
      ref="terminalBody"
    >
      <div 
        v-for="(line, idx) in output" 
        :key="idx"
        :class="{
          'text-green-accent-3': line.type === 'command',
          'text-white': line.type === 'output',
          'text-error': line.type === 'error'
        }"
        class="mb-1"
      >
        {{ line.text }}
      </div>
    </div>
    
    <div class="terminal-input pa-2 border-t d-flex align-center" style="background-color: #1a1a1a;">
      <span class="text-green-accent-3 mr-2 font-mono ml-2">$</span>
      <input 
        v-model="currentCommand" 
        @keyup.enter="executeCommand"
        type="text" 
        class="flex-grow-1 text-white font-mono bg-transparent"
        style="border: none; outline: none; width: 100%;"
        placeholder="Enter command..."
        autocomplete="off"
        spellcheck="false"
      />
    </div>
  </div>
</template>

<style scoped>
.font-mono {
  font-family: 'SFMono-Regular', Consolas, 'Liberation Mono', Menlo, Courier, monospace;
}
.border-t {
  border-top: 1px solid #333;
}
</style>
