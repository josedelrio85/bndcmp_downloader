<script setup>
  import { ref, computed } from 'vue'
    // import HelloWorld from './components/HelloWorld.vue'
  const url = ref('')
  const response = ref('')
  const isLoading = ref(false)

  console.log(import.meta.env.VITE_ENVIRONMENT)
  const apiBaseUrl = import.meta.env.VITE_API_BASE_URL

  const submitRequest = async () => {
    const apiUrl = `${apiBaseUrl}/api/v1/scrapp?url=${encodeURIComponent(url.value)}`

    try {
      isLoading.value = true
      console.log('Requesting:', apiUrl)
      const res = await fetch(apiUrl, {
        method: 'GET',
      })

      if (!res.ok) {
        throw new Error(`HTTP error! status: ${res.status}`)
      }

      const data = await res.text()
      response.value = data
    } catch (error) {
      response.value = `Error: ${error.message}`
    } finally {
      isLoading.value = false
    }
  }
</script>

<template>
  <div class="container">
    <h1>Bandcamp Downloader</h1>
    <div class="form">
      <div class="input-group">
        <input type="text" v-model="url" placeholder="Enter URL">
        <button @click="submitRequest" :disabled="isLoading">Download</button>
      </div>
    </div>
    <div v-if="isLoading" class="loading">Processing request...</div>
    <textarea readonly v-model="response"></textarea>
  </div>
</template>

<style scoped>
.container {
  max-width: 600px;
  margin: 0 auto;
  padding: 20px;
}

.form {
  margin-bottom: 20px;
}

.checkboxes {
  display: flex;
  gap: 10px;
  margin-bottom: 10px;
}

.input-group {
  display: flex;
  gap: 10px;
}

input[type="text"] {
  flex-grow: 1;
  padding: 5px;
}

button {
  padding: 5px 10px;
}

textarea {
  width: 100%;
  height: 200px;
  resize: vertical;
}

.loading {
  margin-top: 10px;
  text-align: center;
  font-style: italic;
  color: #666;
}
</style>
