<template>
    <div>
        <v-navigation-drawer  v-model="drawer"
      bottom
      temporary app dark color="indigo darken-5">
            <v-list dense>
                <v-list-item-group>
                    <v-list-item v-for="(item,i) in items" :key="i" router :to="item.route">
                        <v-list-item-icon>
                            <v-icon v-text="item.icon"></v-icon>
                        </v-list-item-icon>
                        <v-list-item-content v-text="item.text" ></v-list-item-content> 

                    </v-list-item>
                </v-list-item-group>

            </v-list>
        </v-navigation-drawer>

        
        <v-toolbar color="deep-purple" dark >
            <v-app-bar-nav-icon @click.stop="drawer = !drawer" class="white--text"></v-app-bar-nav-icon>
            <v-toolbar-title class="text-uppercase white--text">
                <span class="font weight-light">My</span>
                <span>Orders</span>
            </v-toolbar-title>
            <v-spacer></v-spacer>
            <v-btn  color="grey">
                <span>Sing Out</span>
                <v-icon right>mdi-logout</v-icon>
            </v-btn>
        </v-toolbar>
        
        
    </div>
    
</template>
<script>
export default {
    data(){
        return {
            drawer: false,
            clipped: false,
            items: [
                { text: 'DASHBOARD', icon: 'dashboard',route:'/' },
                { text: 'SEARCH', icon: 'mdi-account', route:'/about' },

            ]
        }
    }
   
}
</script>