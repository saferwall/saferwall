<template>
    <aside class="menu sidebar">
        <ul class="menu-list">
            <li v-for="(item, index) in menu" :class="{'current': $route.name == item.title}">
                <router-link :to="(item.path) ? item.path : ''" 
                            :class="{'is-active': item.children && item.active}"
                            :exact="true" 
                            @click.native="toggleDropdown(index)">
                    <span class="icon is-small"><i :class="item.icon"></i></span>
                    {{item.title}}
                    <span class="icon is-small dropdown-icon" :class="{'active': item.active}" v-if="item.children">
                        <i class="ion-ios-arrow-down"></i>
                    </span>
                </router-link>
                <ul class="dropdown-container" 
                    v-if="item.children"
                    :class="{'active': item.active}"
                    :style="{'height': item.dropdownHeight}">
                    <li v-for="child in item.children" :class="{'current': $route.name == child.title}">
                        <router-link :to="child.path">
                            {{child.title}}
                        </router-link>
                    </li>
                </ul>
            </li>
        </ul>
    </aside>
</template>

<script>
export default {
    data(){
        return{
            menu: [
                {
                    title: 'Summary',
                    path: "/summary/df50dd428c2c0a6c2bffc6720b10d690061f1e3e0d1f5ef2f926942cbf4fc69c",
                    active: false,
                    icon: 'ion-android-settings'
                },
                {
                    title: 'Static analysis',
                    path: false,
                    active: false,
                    icon: 'ion-stats-bars',
                    children: [
                        {title: 'PE', path: '/'},
                        {title: 'Strings', path: '/strings/df50dd428c2c0a6c2bffc6720b10d690061f1e3e0d1f5ef2f926942cbf4fc69c'},
                        {title: 'Antivirus', path: '/antivirus/df50dd428c2c0a6c2bffc6720b10d690061f1e3e0d1f5ef2f926942cbf4fc69c'}
                    ]
                },
                {
                    title: 'Dynamic analysis',
                    path: false,
                    active: false,
                    icon: 'ion-ios-analytics',
                    children: [
                        {title: 'API Calls', path: '/'},
                        {title: 'Network', path: '/'},
                        {title: 'Dropped files', path: '/'},
                        {title: 'Memory dumps', path: '/'},
                    ]
                }
            ]
        }
    },

    mounted(){
        this.menu.map((el) => el.active = (el.children) ? true : false)
        // this.menu[0].active = false
        this.menu.map((el) => el.dropdownHeight = (el.active) ? el.children.length * 36 + 'px' : 0)
    },
    
    methods: {
        toggleDropdown(index){
            if(this.menu[index].children){
                let found = false;
                this.menu.forEach((el, i) => {
                    if(el.active && i == index) {
                        el.active = false; 
                        el.dropdownHeight = 0;
                        found = true
                    }
                })
                if(found) return
                // this.menu.map((el) => {el.active = false; el.dropdownHeight = 0})
                this.menu[index].active = true
                this.menu[index].dropdownHeight = this.menu[index].children.length * 36 + 'px'
            }
        }
    }
}
</script>

<style scoped lang="scss">
@import '../../assets/scss/variables';

aside.sidebar{
    background-color: #fff;
    position: fixed;
    top:50px;
    left:0;
    height:calc(100% - 50px);
    width:200px;
    box-shadow: 0 0 30px rgba(black, .05);
    padding-top: 20px;

    .menu-list{
        .is-active{
            background-color: transparent;
            color:#4a4a4a
        }

        .dropdown-icon{
            float:right;

            i{transition: all .2s}

            &.active{
                i{
                    transform:rotate(180deg)
                }
            }
        }

        .current{
            background-color: $primary-color;
            color: #fff;

            a{
                color:#fff!important;
            }
        }

        ul.dropdown-container{
            height: 0;
            overflow:hidden;
            transition: all .2s;
            padding-left:0;
            margin:0 0 0 0.75em;

            &.active{
                height: auto;
                margin: 0 0 0.75em 0.75em;
            }
        }

        li:not(.current) > a:hover{
            background-color: #fff;
            color: $primary-color;
        }
        
        li.current > a:hover{
            background-color: $primary-color;
        }
    }
}
</style>

