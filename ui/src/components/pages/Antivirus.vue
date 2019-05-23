<template>
    <div class="tile is-ancestor">
        <loader v-if="showLoader"></loader>
        <div class="tile is-parent is-6">
            <div class="tile is-child box" v-if="!showLoader">
                <h4 class="title">First Scan</h4>
                <table class="table is-striped is-fullwidth">
                    <thead>
                        <tr>
                            <th>Antivirus</th>
                            <th>Output</th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr v-for="(scan, index, i) of firstScan">
                            <td>{{scan.vendor}}</td>
                            <td>
                                <span :class="[
                                        {'has-text-success': !scan.detected}, 
                                        {'has-text-danger': scan.detected}
                                    ]"
                                    style="position:relative"
                                    @mouseover="mouseOver('first', index)"
                                    @mouseleave="mouseLeave('first', index)">
                                    <span :class="{'transparent': scan.detected && scan.showCopy}">
                                        <i class="output-icon icon" :class="[{'ion-alert-circled': scan.detected}, {'ion-checkmark-circled': !scan.detected}]"></i>
                                        {{scan.output}}
                                    </span>
                                    <transition name="fade">
                                        <copy v-if="scan.detected && scan.showCopy" :content="scan.output"></copy>
                                    </transition> 
                                </span>
                            </td>
                        </tr>
                    </tbody>
                </table>
            </div>
        </div>
        <div class="tile is-parent is-6">
            <div class="tile is-child box" v-if="!showLoader">
                <h4 class="title">Last Scan</h4>
                <table class="table is-striped is-fullwidth">
                    <thead>
                        <tr>
                            <th>Antivirus</th>
                            <th>Output</th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr v-for="(scan, index) of lastScan">
                            <td>{{scan.vendor}}</td>
                            <td>
                                <span :class="[
                                        {'has-text-success': !scan.detected}, 
                                        {'has-text-danger': scan.detected}
                                    ]"
                                    style="position:relative"
                                    @mouseover="mouseOver('last', index)"
                                    @mouseleave="mouseLeave('last', index)">
                                    <span :class="{'transparent': scan.detected && scan.showCopy}">
                                        <i class="output-icon icon" :class="[{'ion-alert-circled': scan.detected}, {'ion-checkmark-circled': !scan.detected}]"></i>
                                        {{scan.output}}
                                    </span>
                                    <transition name="fade">
                                        <copy v-if="scan.detected && scan.showCopy" :content="scan.output"></copy>
                                    </transition> 
                                </span>
                            </td>
                        </tr>
                    </tbody>
                </table>
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
        return {
            showLoader: true,
            lastScan: {},
            firstScan: {}
        }
    },
    methods: {
        mouseOver(type, index){
            if(type == 'first'){
                var temp = []
                this.firstScan.forEach((e, i) => {
                    let output = e.output
                    let vendor = e.vendor
                    let detected = e.detected
                    let showCopy = ( i == index ) ? true : e.showCopy
                    temp.push({output, vendor, detected, showCopy})
                })
                this.firstScan = []
                this.firstScan = temp
            }else{
                var temp = []
                this.lastScan.forEach((e, i) => {
                    let output = e.output
                    let vendor = e.vendor
                    let detected = e.detected
                    let showCopy = ( i == index ) ? true : e.showCopy
                    temp.push({output, vendor, detected, showCopy})
                })
                this.lastScan = []
                this.lastScan = temp
            }
        },
        mouseLeave(type, index){
            if(type == 'first'){
                var temp = []
                this.firstScan.forEach((e, i) => {
                    let output = e.output
                    let vendor = e.vendor
                    let detected = e.detected
                    let showCopy = ( i == index ) ? false : e.showCopy
                    temp.push({output, vendor, detected, showCopy})
                })
                this.firstScan = []
                this.firstScan = temp
            }else{
                var temp = []
                this.lastScan.forEach((e, i) => {
                    let output = e.output
                    let vendor = e.vendor
                    let detected = e.detected
                    let showCopy = ( i == index ) ? false : e.showCopy
                    temp.push({output, vendor, detected, showCopy})
                })
                this.lastScan = []
                this.lastScan = temp
            }
        }
    },
    mounted(){
        let url = Global.apiUrl + this.$route.params.hash + '?api-key=' + Global.apiKey
        axios.get(url)
            .then((data) => {
                this.showLoader = false
                this.firstScan = data.data.multiav
                this.lastScan = data.data.multiav

								for (let key in this.firstScan) {
										const first = this.firstScan[key];
                    first.showCopy = false
                    const last = this.lastScan[key];
                    last.showCopy = false
								}
            })
            .catch(err => console.error(err))
    }

}
</script>
<style lang="scss" scoped>
.fade-enter-active, .fade-leave-active {
  transition-property: opacity;
  transition-duration: .25s;
}

.fade-enter-active {
  transition-delay: 0;
}

.fade-enter, .fade-leave-active {
  opacity: 0
}
span{transition: all .2s}
.transparent{opacity: .35}
.output-icon{
    font-size: 18px;
    display:inline-block;
}
</style>

