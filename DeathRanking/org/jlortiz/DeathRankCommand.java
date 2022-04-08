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

import java.util.ArrayList;
import java.util.Comparator;
import java.util.Map;
import java.util.UUID;
import org.bukkit.Bukkit;
import org.bukkit.ChatColor;
import org.bukkit.command.Command;
import org.bukkit.command.CommandExecutor;
import org.bukkit.command.CommandSender;

/**
 *
 * @author jlortiz
 */
public class DeathRankCommand implements CommandExecutor {
    
    private final DeathRanking plugin;
    private Map<String, Object> deaths;
    
    public DeathRankCommand(DeathRanking plugin) {
        this.plugin = plugin;
        this.deaths = this.plugin.config.getValues(false);
    }
    
    @Override
    public boolean onCommand(CommandSender sender, Command cmd, String alias, String[] args) {
        if (this.plugin.hasChanged) {
            this.deaths = this.plugin.config.getValues(false);
        }
        int page = 1;
        if (args.length > 0) {
            try {
                page = Integer.parseInt(args[0]);
            } catch (NumberFormatException e) {
                sender.sendMessage(ChatColor.RED+args[0]+" is not a number!");
                return true;
            }
        }
        if (page < 1)
            page = 1;
        if (page > Math.ceil(this.deaths.size()/8.0))
            page = (int)Math.ceil(this.deaths.size()/8.0);
        StringBuilder message = new StringBuilder(ChatColor.YELLOW.toString());
        message.append(" ---- ").append(ChatColor.GOLD).append("Death Rankings").append(ChatColor.YELLOW).append(" -- ");
        message.append(ChatColor.GOLD).append("Page ").append(ChatColor.RED).append(page).append("/").append((int)Math.ceil(deaths.size()/8.0));
        message.append(ChatColor.YELLOW).append(" ----\n").append(ChatColor.GOLD).append("Total server deaths: ");
        if (deaths.isEmpty()) {
            sender.sendMessage(ChatColor.GOLD+"Surprisingly, nobody has died yet.");
            return true;
        }
        int total = 0;
        for (Object i: this.deaths.values()) {
            total += (Integer)i;
        }
        message.append(ChatColor.RED).append(total).append("\n").append(ChatColor.WHITE);
        ArrayList<PlayerDeathCount> dList = new ArrayList<>(this.deaths.size());
        for (String k: this.deaths.keySet()) {
            dList.add(new PlayerDeathCount(this.plugin.config.getInt(k), Bukkit.getOfflinePlayer(UUID.fromString(k)).getName()));
        }
        dList.sort(new PlayerDeathSorter());
        for (int i=(page-1)*8; i<page*8; i++) {
            if (i>=dList.size())
                break;
            message.append(i+1).append(". ").append(dList.get(i).player).append(", ");
            message.append(dList.get(i).deaths).append("\n");
        }
        message.deleteCharAt(message.length()-1);
        sender.sendMessage(message.toString());
        return true;
    }
    
    private class PlayerDeathCount {
        public final int deaths;
        public final String player;
        
        public PlayerDeathCount(int deaths, String player) {
            this.deaths = deaths;
            this.player = player;
        }
    }
    
    private class PlayerDeathSorter implements Comparator<PlayerDeathCount> {
        @Override
        public int compare(PlayerDeathCount t1, PlayerDeathCount t2) {
            return t2.deaths - t1.deaths;
        }
    }
}
