<template>
    <div class="svg-icon">
        <svg width="50px"
            height="50px"  
            xmlns="http://www.w3.org/2000/svg" 
            viewBox="0 0 100 100" 
            preserveAspectRatio="xMidYMid" 
            class="lds-dual-ring" 
            style="background: none;">
            
            <circle cx="50" 
                    cy="50" 
                    
                    fill="none" 
                    stroke-linecap="round" 
                    r="40" 
                    stroke-width="4" 
                    stroke="#363636" 
                    stroke-dasharray="62.83185307179586 62.83185307179586" 
                    transform="rotate(252 50 50)">
            
                    <animateTransform attributeName="transform" 
                                        type="rotate" 
                                        calcMode="linear" 
                                        values="0 50 50;360 50 50" 
                                        keyTimes="0;1" 
                                        dur="1s" 
                                        begin="0s" 
                                        repeatCount="indefinite">
                    </animateTransform>
            </circle>
        </svg>
    </div>
</template>
<style lang="scss" scoped>
.svg-icon{
    position: absolute;top:70px;left:50%;
    text-align:center;
    transform:translateX(-50%);
}
</style>

