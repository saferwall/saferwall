<template>
    <div>
        <loader v-if="showLoader"></loader>
        <div class="tile is-ancestor" v-if="!showLoader">
            <div class="tile is-parent">
                <div class="tile is-child box">
                    <h4 class="title">Basic Properties</h4>
                    <div v-for="(i, index) in summaryData" 
                        v-if="index !== 'av'" 
                        class="data-data">
                        <strong class="data-label">
                            {{
                                (uppercaseFields.includes(index)) ? index.toUpperCase() : 
                                (index == 'filesize' ? 'File Size' : 
                                (index == 'trid' ? 'TRiD' : 
                                (index == 'ssdeep' ? 'SSDeep' : index )))
                            }}
                        </strong>
                        <span class="data-value" v-if="index !== 'trid'">
                            <span class="value-text">{{(index !== 'sha-512') ? i : i.substring(0, 70) + '...'}}</span>

                            <copy :content="i"></copy>
                        </span>
                        <span class="data-value" :class="{'trid-container': index == 'trid'}" v-if="index == 'trid'">
                            <p v-for="t in summaryData.trid">
                            <span class="trid">
                                <span class="value-text">{{t}}</span>

                                <copy :content="t"></copy>
                            </span>
                            </p>
                        </span>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>
<script>
import axios from 'axios'
import Global from '@/global'
import Loader from '@/components/elements/Loader'
import Copy from '@/components/elements/Copy'

export default {
    components: {
        'loader': Loader,
        'copy': Copy
    },
    data(){
        return{
            showLoader: true,
            summaryData: {},
            uppercaseFields: ['md5', 'sha-1', 'sha-256', 'sha-512']
        }
    },
    methods: {
        bytesToSize(bytes) {
            var sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB'];
            if (bytes == 0) return '0 Byte';
            var i = parseInt(Math.floor(Math.log(bytes) / Math.log(1024)));
            return Math.round(bytes / Math.pow(1024, i), 2) + ' ' + sizes[i];
        }
    },
    mounted(){
      axios.get(Global.apiUrl + this.$route.params.hash + '?api-key=' + Global.apiKey)
        .then(data => {
            this.showLoader = false

            data.data['sha-1'] = data.data.sha1
            data.data['sha-256'] = data.data.sha256
            data.data['sha-512'] = data.data.sha512
            delete data.data.sha1
            delete data.data.sha256
            delete data.data.sha512
            
            this.summaryData.filesize = this.bytesToSize(data.data.filesize)
            this.summaryData.magic = data.data.magic
            this.summaryData.md5 = data.data.md5
            this.summaryData['sha-1'] = data.data['sha-1']
            this.summaryData['sha-256'] = data.data['sha-256']
            this.summaryData['sha-512'] = data.data['sha-512']
            this.summaryData.ssdeep = data.data.ssdeep
            this.summaryData.trid = data.data.trid
            console.log(this.summaryData)
        })
    }
}
</script>
<style lang="scss" scoped>
.data-data{
    float: left;
    width:100%;
    padding: 5px;

    &:nth-child(even){
        background: rgba(black, .03)
    }

    &:not(:last-child), .trid:not(:last-child) {margin-bottom: 3px;}
    
    .data-label{
        float:left;
        width:70px;
        text-transform: capitalize;
    }

    .data-value{
        float:left;

        .value-text{transition: opacity .2s;}

        .copy{
            opacity: 0;
            transition: opacity .2s;
        }

        &:not(.trid-container):hover{
            .value-text{opacity: .35;}
            & > .copy{opacity: 1}
        }
    }

    .trid, .data-value{
        position:relative;
    }

    .trid{
        position: relative;

        &:hover{
            .value-text{opacity: .35}
            & > .copy{opacity: 1}
        }
    }
}
</style>

