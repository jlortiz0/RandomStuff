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

import org.bukkit.Bukkit;
import org.bukkit.ChatColor;
import org.bukkit.command.Command;
import org.bukkit.command.CommandExecutor;
import org.bukkit.command.CommandSender;
import org.bukkit.entity.Player;

/**
 *
 * @author jlortiz
 */
public class DeathsCommand implements CommandExecutor {
    
    private final DeathRanking plugin;
    
    public DeathsCommand(DeathRanking plugin) {
        this.plugin = plugin;
    }
    
    @Override
    public boolean onCommand(CommandSender sender, Command cmd, String alias, String[] args) {
        if (!(sender instanceof Player)) {
            sender.sendMessage(ChatColor.RED+"This command is for players only.");
            return true;
        }
        if (args.length > 0) {
            String uuid = Bukkit.getOfflinePlayer(args[0]).getUniqueId().toString();
            if (!this.plugin.config.contains(uuid)) {
                sender.sendMessage(ChatColor.BLUE+"You haven't died yet. Impressive.");
            } else {
                sender.sendMessage(ChatColor.BLUE+"You have died "+this.plugin.config.getInt(uuid)+" times.");
            }
        } else {
        String uuid = ((Player)sender).getUniqueId().toString();
            if (!this.plugin.config.contains(uuid)) {
                sender.sendMessage(ChatColor.BLUE+"You haven't died yet. Impressive.");
            } else {
                sender.sendMessage(ChatColor.BLUE+"You have died "+this.plugin.config.getInt(uuid)+" times.");
            }
        }
        return true;
    }
}
