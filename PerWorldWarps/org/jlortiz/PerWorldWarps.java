/*
 * Copyright (C) 2019 jlortiz
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */
package org.jlortiz;

import java.io.File;
import java.io.FileWriter;
import java.io.IOException;
import java.util.HashMap;
import java.util.Map;
import java.util.logging.Level;
import java.util.logging.Logger;
import org.bukkit.ChatColor;
import org.bukkit.Location;
import org.bukkit.World;
import org.bukkit.command.Command;
import org.bukkit.command.CommandSender;
import org.bukkit.command.TabExecutor;
import org.bukkit.entity.LivingEntity;
import org.bukkit.event.Listener;
import org.bukkit.event.world.WorldLoadEvent;
import org.bukkit.plugin.java.JavaPlugin;
import org.yaml.snakeyaml.Yaml;

/**
 *
 * @author jlortiz
 */
public class PerWorldWarps extends JavaPlugin implements Listener,TabExecutor {
    
    @Override
    public void onEnable() {
        this.getDataFolder().mkdir();
        for (World w: getServer().getWorlds()) {
            new File(this.getDataFolder(), w.getName()).mkdir();
        }
        this.getServer().getPluginManager().registerEvents(this, this);
        this.getCommand("warp").setExecutor(new WarpCommand(this));
        this.getCommand("setWarp").setExecutor(this);
        this.getCommand("delWarp").setExecutor(this);
    }
    
    public void onWorldLoad(WorldLoadEvent event) {
        new File(this.getDataFolder(), event.getWorld().getName()).mkdir();
    }
    
    @Override
    public boolean onCommand(CommandSender sender, Command cmd, String alias, String[] args) {
        if (!(sender instanceof LivingEntity)) {
            sender.sendMessage("This command for players only.");
            return true;
        }
        if (args.length==0) {
            return false;
        }
        LivingEntity plr = (LivingEntity)sender;
        File warp = new File(this.getDataFolder(), plr.getWorld().getName()+"/"+args[0]+".yml");
        if (cmd.getName().equals("setwarp")) {
            Map<String, Object> warpMap = new HashMap<>();
            Location loc = plr.getLocation();
            warpMap.put("x", loc.getX());
            warpMap.put("y", loc.getY());
            warpMap.put("z", loc.getZ());
            warpMap.put("pitch", plr.getEyeLocation().getPitch());
            warpMap.put("yaw", plr.getEyeLocation().getYaw());
            try {
                warp.createNewFile();
                new Yaml().dump(warpMap, new FileWriter(warp));
                sender.sendMessage(ChatColor.GOLD+"Warp "+ChatColor.RED+args[0]+ChatColor.GOLD+" created.");
            } catch (IOException ex) {
                Logger.getLogger("Minecraft").log(Level.WARNING, "Error creating warp", ex);
                sender.sendMessage(ChatColor.RED+"Error: "+ChatColor.DARK_RED+" Failed to create warp. Check console for details.");
            }
        } else {
            if (!warp.exists()) {
                sender.sendMessage(ChatColor.RED+"Error: "+ChatColor.DARK_RED+"That warp does not exist.");
                return true;
            }
            warp.delete();
            sender.sendMessage(ChatColor.GOLD+"Warp "+ChatColor.RED+args[0]+ChatColor.GOLD+" has been removed.");
        }
        return true;
    }
}
