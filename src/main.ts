import { createApp } from 'vue'
import './style.css'
import App from './App.vue'

// Vuetify
import 'vuetify/styles'
import '@mdi/font/css/materialdesignicons.css'
import { createVuetify } from 'vuetify'
import * as components from 'vuetify/components'
import * as directives from 'vuetify/directives'

const vuetify = createVuetify({
  components,
  directives,
  theme: {
    defaultTheme: 'dark' // Twitter-like dark mode as default
  }
})

// Initialize Go WASM
const go = new (window as any).Go()
WebAssembly.instantiateStreaming(fetch('/main.wasm'), go.importObject).then((result) => {
  go.run(result.instance)
})

createApp(App).use(vuetify).mount('#app')
