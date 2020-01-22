<template>
  <aside class="menu sidebar">
    <ul class="menu-list">
      <li
        v-for="(item, index) in menuItems"
        :class="{ current: route.toLowerCase() === item.title.toLowerCase() }"
        :key="index"
      >
        <router-link
          :to="item.path ? item.path : ''"
          :class="{ 'is-active': item.children && item.active }"
          :exact="true"
          @click.native="toggleDropdown(index)"
        >
          <span class="icon is-small"><i :class="item.icon"></i></span>
          {{ item.title }}
          <span
            class="icon is-small dropdown-icon"
            :class="{ active: item.active }"
            v-if="item.children"
          >
            <i class="ion-ios-arrow-down"></i>
          </span>
        </router-link>
        <ul
          class="dropdown-container"
          v-if="item.children"
          :class="{ active: item.active }"
          :style="{ height: item.dropdownHeight }"
        >
          <li
            v-for="(child, index) in item.children"
            :class="{ current: route.toLowerCase() === child.title.toLowerCase() }"
            :key="index"
          >
            <router-link :to="child.path">
              {{ child.title }}
            </router-link>
          </li>
        </ul>
      </li>
    </ul>
  </aside>
</template>

<script>
import {mapGetters} from 'vuex'
export default {
  data() {
    return {
      menu: [
        {
          title: "Summary",
          slug: "summary",
          active: false,
          icon: "ion-android-settings",
        },
        {
          title: "Static analysis",
          path: false,
          active: false,
          icon: "ion-stats-bars",
          children: [
            // { title: "PE", slug: "pe" },
            {
              title: "Strings",
              slug: "strings",
            },
            {
              title: "Antivirus",
              slug: "antivirus",
            },
          ],
        },
        // {
        //   title: "Dynamic analysis",
        //   path: false,
        //   active: false,
        //   icon: "ion-ios-analytics",
        //   children: [
        //     { title: "API Calls", path: "/" },
        //     { title: "Network", path: "/" },
        //     { title: "Dropped files", path: "/" },
        //     { title: "Memory dumps", path: "/" },
        //   ],
        // },
      ],
    }
  },
  computed: {
    ...mapGetters(['getHashContext']),
    menuItems: function() {
      const hash  = this.getHashContext
      return this.menu.map(({ slug, children, ...item }) => ({
        ...item,
        ...(slug ? { path: `/${slug}/${hash}` } : {}),
        ...(children
          ? {
              children: children.map(({ slug, ...item }) => ({
                ...item,
                ...(slug ? { path: `/${slug}/${hash}` } : {}),
              })),
            }
          : {}),
      }))
    },
    route: function(){
      return this.$route.path.replace(/[/]/g, '')
    }
  },

  mounted() {
    this.menu.map((el) => {
      el.active = Boolean(el.children)
    })
    // this.menu[0].active = false
    this.menu.map((el) => {
      el.dropdownHeight = el.active ? el.children.length * 36 + "px" : 0
    })
  },

  methods: {
    toggleDropdown(index) {
      if (this.menu[index].children) {
        let found = false
        this.menu.forEach((el, i) => {
          if (el.active && i === index) {
            el.active = false
            el.dropdownHeight = 0
            found = true
          }
        })
        if (found) return
        // this.menu.map((el) => {el.active = false; el.dropdownHeight = 0})
        this.menu[index].active = true
        this.menu[index].dropdownHeight =
          this.menu[index].children.length * 36 + "px"
      }
    },
  },
}
</script>

<style scoped lang="scss">
@import "../../assets/scss/variables";
$header-height: 50px;
aside.sidebar {
  background-color: #fff;
  position: fixed;
  top: $header-height;
  left: 0;
  height: calc(100% - #{ $header-height });
  width: 200px;
  box-shadow: 0 0 30px rgba(black, 0.05);
  padding-top: 20px;

  .menu-list {
    .is-active {
      background-color: transparent;
      color: #4a4a4a;
    }

    .dropdown-icon {
      float: right;

      i {
        transition: all 0.2s;
      }

      &.active {
        i {
          transform: rotate(180deg);
        }
      }
    }

    .current {
      background-color: $primary-color;
      color: #fff;

      a {
        color: #fff !important;
      }
    }

    ul.dropdown-container {
      height: 0;
      overflow: hidden;
      transition: all 0.2s;
      padding-left: 0;
      margin: 0 0 0 0.75em;

      &.active {
        height: auto;
        margin: 0 0 0.75em 0.75em;
      }
    }

    li:not(.current) > a:hover {
      background-color: #fff;
      color: $primary-color;
    }

    li.current > a:hover {
      background-color: $primary-color;
    }
  }
}
</style>
