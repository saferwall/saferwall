<template>
  <section class="main-content" :class="{'fullwidth': fullwidth}">
    <div class="container is-fluid">
        <nav class="breadcrumb" aria-label="breadcrumbs" v-if="!fullwidth">
            <ul>
                <li><router-link to="/">Dashboard</router-link></li>
                <li class="is-active"><a href="#" aria-current="page">{{$route.name}}</a></li>
            </ul>
        </nav>
        <slot></slot>
    </div>
  </section>
</template>
<script>
export default {
    props: ['fullwidth']
}
</script>
<style scoped lang="scss">
@import "../../assets/scss/variables";
section.main-content{
    float:right;
    padding-top: 20px;
    margin-top:50px;

    &:not(.fullwidth){width:calc(100% - 200px);}

    &.fullwidth{width:100%;}
}
.breadcrumb{
    a{
        color: $primary-color;
    }
}
</style>

