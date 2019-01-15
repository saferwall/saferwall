<template>
  <div class="copy" :class="{'copied': copied}" @click="copy($event)">
      Copy
  </div>
</template>
<script>
export default {
    props: ['content'],
    data(){
        return{
            copied: false
        }
    },
    methods: {
        copy(e){
            this.copied = true
            Clipboard.copy(this.content);
            setTimeout(() => {this.copied = false}, 500)
        }
    },
    created(){
        window.Clipboard = (function(window, document, navigator) {
            var textArea,
                copy;

            function isOS() {
                return navigator.userAgent.match(/ipad|iphone/i);
            }

            function createTextArea(text) {
                textArea = document.createElement('textArea');
                textArea.value = text;
                document.body.appendChild(textArea);
            }

            function selectText() {
                var range,
                    selection;

                if (isOS()) {
                    range = document.createRange();
                    range.selectNodeContents(textArea);
                    selection = window.getSelection();
                    selection.removeAllRanges();
                    selection.addRange(range);
                    textArea.setSelectionRange(0, 999999);
                } else {
                    textArea.select();
                }
            }

            function copyToClipboard() {        
                document.execCommand('copy');
                document.body.removeChild(textArea);
            }

            copy = function(text) {
                createTextArea(text);
                selectText();
                copyToClipboard();
            };

            return {
                copy: copy
            };
        })(window, document, navigator);
    }
}
</script>
<style lang="scss" scoped>
@import '../../assets/scss/variables';
@keyframes copying{
    from{
        transform:translate(50%, -50%) scale(1);
        opacity:1;
    }
    to{
        transform:translate(50%, -50%) scale(1.2);
        opacity:0;
    }
}
.copy{
    position:absolute;
    top:50%;
    right: 50%;
    transform:translate(50%, -50%);
    display:inline-block;
    color:#fff;
    font-weight:500;
    font-size:12px;
    font-weight: 500;
    padding: 0 3px;
    border-radius:3px;
    cursor:pointer;
    background-color: $primary-color;
    z-index: 999;

    &.copied{
        animation: copying .5s
    }
}
</style>

