<template>
  <div class="app-container">
    <div class="filter-container">
      <div class="box-card">
        <video id="videoElement" controls="true" height="100%" width="100%" />
      </div>
    </div>
  </div>
</template>

<script>

import flvPlayer from 'flv.js'
import { fetchArticle } from '@/api/article'
const stream = {
}
export default {
  name: 'Play',
  data() {
    return {
      flvPlayer: null,
      stream: Object.assign({}, stream)
    }
  },
  created() {
    const id = this.$route.params && this.$route.params.id
    this.fetchData(id)
  },
  mounted() {
  },
  methods: {
    fetchData(id) {
      fetchArticle(id).then(response => {
        this.stream = response.data
        setTimeout(() => {
          if (flvPlayer.isSupported()) {
            var videoElement = document.getElementById('videoElement')
            this.flvPlayer = flvPlayer.createPlayer({
              url: this.stream.FlvUrl,
              type: 'flv',
              isLive: true
            })
            this.flvPlayer.attachMediaElement(videoElement)
            this.flvPlayer.load()
          }
        }, 2 * 1000)
      }).catch(err => {
        console.log(err)
      })
    }
  }
}
</script>
