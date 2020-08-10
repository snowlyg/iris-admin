<template>
  <div class="app-container">
    <div class="filter-container">
      <div class="box-card">
        <div id="app">
          <video id="hlsVideo" ref="hlsVideo" controls preload="true" />
        </div>
      </div>
    </div>
  </div>
</template>
<style lang="less" scoped>
  #hlsVideo {
    width: 100%;
    height: 100%;
    border: none;
  }
</style>
<script>
import Hls from 'hls.js'
import { fetchArticle } from '@/api/article'
const stream = {}
export default {
  name: 'Play',
  data() {
    return {
      hls: '',
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
        var video = document.getElementById('hlsVideo')
        if (Hls.isSupported()) {
          var hls = new Hls()
          hls.loadSource(this.stream.HlsUrl)
          hls.attachMedia(video)
          hls.on(Hls.Events.MANIFEST_PARSED, function() {
            video.play()
          })
        } else if (video.canPlayType('application/vnd.apple.mpegurl')) {
          video.src = this.stream.HlsUrl
          video.addEventListener('loadedmetadata', function() {
            video.play()
          })
        }
      }).catch(err => {
        console.log(err)
      })
    }

  }
}
</script>
