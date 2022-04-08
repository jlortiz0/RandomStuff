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
import java.io.IOException;
import java.util.logging.Level;
import java.util.logging.Logger;
import org.bukkit.Bukkit;
import org.bukkit.ChatColor;
import org.bukkit.configuration.InvalidConfigurationException;
import org.bukkit.configuration.file.FileConfiguration;
import org.bukkit.entity.Player;
import org.bukkit.event.EventHandler;
import org.bukkit.event.Listener;
import org.bukkit.event.entity.PlayerDeathEvent;
import org.bukkit.plugin.java.JavaPlugin;

/**
 *
 * @author jlortiz
 */
public class DeathRanking extends JavaPlugin implements Listener {

    public final FileConfiguration config = this.getConfig();
    public boolean hasChanged = false;
    
    /**
     * @param args the command line arguments
     */
    @Override
    public void onEnable() {
        try {
            this.getDataFolder().mkdir();
            new File(this.getDataFolder(), "deaths.yml").createNewFile();
        }
        catch (IOException ex) {
            Logger.getLogger(DeathRanking.class.getName()).log(Level.SEVERE, "Could not create deaths.yml!", ex);
            this.setEnabled(false);
            return;
        }
        try {
            config.load(new File(this.getDataFolder(), "deaths.yml"));
        }
        catch (IOException | InvalidConfigurationException ex) {
            Logger.getLogger(DeathRanking.class.getName()).log(Level.SEVERE, "Could not load deaths.yml!", ex);
            this.setEnabled(false);
            return;
        }
        this.getServer().getPluginManager().registerEvents(this, this);
        this.getCommand("deaths").setExecutor(new DeathsCommand(this));
        this.getCommand("deathtop").setExecutor(new DeathRankCommand(this));
    }
    
    @Override
    public void onDisable() {
        try {
            config.save(new File(this.getDataFolder(), "deaths.yml"));
        }
        catch (IOException ex) {
            Logger.getLogger(DeathRanking.class.getName()).log(Level.SEVERE, "Failed to save deaths!", ex);
        }
    }
    
    @EventHandler(ignoreCancelled=true)
    public void onPlayerDeath(PlayerDeathEvent event) {
        String uuid = event.getEntity().getUniqueId().toString();
        if (!config.contains(uuid)) {
            config.set(uuid, 0);
        }
        int deaths = config.getInt(uuid);
        deaths++;
        StringBuilder message = new StringBuilder(ChatColor.BLUE.toString());
        message.append("This is your ").append(deaths);
        switch (deaths%10) {
            case 1:
                message.append("st");
                break;
            case 2:
                message.append("nd");
                break;
            case 3:
                message.append("rd");
                break;
            default:
                message.append("th");
        }
        message.append(" death.");
        if (deaths%50 == 0) {
            message.append(ChatColor.GREEN).append(" Congratulations!");
            for (Player player : Bukkit.getOnlinePlayers()) {
                player.sendMessage(ChatColor.GREEN+event.getEntity().getName()+" has died "+deaths+" times!");
            }
        }
        event.getEntity().sendMessage(message.toString());
        config.set(uuid, deaths);
        this.hasChanged = true;
        try {
            config.save(new File(this.getDataFolder(), "deaths.yml"));
        }
        catch (IOException ex) {
            Logger.getLogger(DeathRanking.class.getName()).log(Level.SEVERE, "Failed to save deaths!", ex);
        }
    }
}
