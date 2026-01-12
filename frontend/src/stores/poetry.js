import { defineStore } from 'pinia'
import { ref } from 'vue'

export const usePoetryStore = defineStore('poetry', () => {
  const dynasties = ref([])
  const categories = ref([])
  const currentPoem = ref(null)
  const loading = ref(false)

  function setDynasties(data) {
    dynasties.value = data
  }

  function setCategories(data) {
    categories.value = data
  }

  function setCurrentPoem(poem) {
    currentPoem.value = poem
  }

  return {
    dynasties,
    categories,
    currentPoem,
    loading,
    setDynasties,
    setCategories,
    setCurrentPoem
  }
})
