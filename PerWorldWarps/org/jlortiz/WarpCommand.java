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

import java.io.*;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.List;
import java.util.Map;
import java.util.logging.Level;
import java.util.logging.Logger;
import org.bukkit.ChatColor;
import org.bukkit.Location;
import org.bukkit.command.Command;
import org.bukkit.command.CommandSender;
import org.bukkit.command.TabExecutor;
import org.bukkit.entity.LivingEntity;
import org.yaml.snakeyaml.Yaml;
import org.yaml.snakeyaml.error.YAMLException;

/**
 *
 * @author jlortiz
 */
public class WarpCommand implements TabExecutor {
    
    private final PerWorldWarps plugin;
    private static final Yaml yaml = new Yaml();
    private static final Logger log = Logger.getLogger("Minecraft");
    
    public WarpCommand(PerWorldWarps plugin) {
        this.plugin = plugin;
    }
    @Override
    public boolean onCommand(CommandSender sender, Command cmd, String alias, String[] args) {
        if (!(sender instanceof LivingEntity)) {
            sender.sendMessage("This command for players only.");
            return true;
        }
        if (args.length==0) {
            onCommand(sender, cmd, alias, new String[] {"1"});
        } else {
            LivingEntity plr = (LivingEntity)sender;
            try {
                int page = Integer.parseInt(args[0]);
                if (page<1)
                    page=1;
                StringBuilder message = new StringBuilder(ChatColor.GOLD.toString());
                message.append("There are ");
                String[] warps = new File(plugin.getDataFolder(), plr.getWorld().getName()).list();
                if (warps==null || warps.length==0) {
                    sender.sendMessage(ChatColor.GOLD+"There are no warps in "+ChatColor.RED+plr.getWorld().getName());
                    return true;
                }
                Arrays.sort(warps);
                message.append(ChatColor.RED);
                message.append(warps.length);
                message.append(ChatColor.GOLD);
                message.append(" warps in ");
                message.append(ChatColor.RED);
                message.append(plr.getWorld().getName());
                message.append(ChatColor.GOLD);
                message.append(".");
                if (warps.length>20) {
                    if (page>Math.ceil(warps.length/20.0))
                        page=(int)Math.ceil(warps.length/20.0);
                    message.append(" Showing page ");
                    message.append(ChatColor.RED);
                    message.append(page);
                    message.append(ChatColor.GOLD);
                    message.append(" of ");
                    message.append(ChatColor.RED);
                    message.append((int)Math.ceil(warps.length/20.0));
                }
                message.append("\n");
                message.append(ChatColor.RESET);
                for (int i=(page-1)*20; i<page*20; i++) {
                    if (i>=warps.length)
                        break;
                    String s = warps[i];
                    message.append(s.substring(0, s.length()-4));
                    message.append(" ");
                }
                message.deleteCharAt(message.length()-1);
                sender.sendMessage(message.toString());
            } catch (NumberFormatException e) {
                if (args[0].equals("importFromEss")) {
                    int errorCount = 0;
                    for (File f: new File(plugin.getDataFolder(), "../Essentials/warps").listFiles()) {
                        Map<String, Object> warpMap;
                        try {
                            warpMap = yaml.load(new FileInputStream(f));
                        } catch (FileNotFoundException ex) {
                            errorCount++;
                            log.log(Level.WARNING, "Error importing warps from Essentials", ex);
                            continue;
                        }
                        String world = (String)warpMap.get("world");
                        warpMap.remove("world");
                        String name = (String)warpMap.get("name");
                        warpMap.remove("name");
                        warpMap.remove("lastowner");
                        try {
                            new File(plugin.getDataFolder(), world+"/"+name+".yml").createNewFile();
                            yaml.dump(warpMap, new FileWriter(new File(plugin.getDataFolder(), world+"/"+name+".yml")));
                        } catch (IOException ex) {
                            errorCount++;
                            log.log(Level.WARNING, "Error importing warps from Essentials", ex);
                        }
                    }
                    if (errorCount==0) {
                        sender.sendMessage(ChatColor.GREEN+"Success! Imported warps without issue.");
                    } else {
                        sender.sendMessage(ChatColor.GOLD+"There were "+ChatColor.RED+errorCount+ChatColor.GOLD+" errors importing warps from Essentials. Check console for details.");
                    }
                } else {
                    File warp = new File(plugin.getDataFolder(), plr.getWorld().getName()+"/"+args[0]+".yml");
                    if (!warp.exists()) {
                        sender.sendMessage(ChatColor.RED+"Error: "+ChatColor.DARK_RED+"That warp does not exist.");
                        return true;
                    }
                    Map<String, Object> warpMap;
                    try {
                        warpMap = yaml.load(new FileInputStream(warp));
                    } catch (FileNotFoundException ex) {
                        log.log(Level.WARNING, "Error warping to "+args[0], ex);
                        sender.sendMessage(ChatColor.RED+"Error: "+ChatColor.DARK_RED+"That warp does not exist.");
                        return true;
                    } catch (YAMLException ex) {
                        log.log(Level.WARNING, "Error warping to "+args[0], ex);
                        sender.sendMessage(ChatColor.RED+"Error: "+ChatColor.DARK_RED+"That warp is corrupt. Please inform the admins.");
                        return true;
                    }
                    try {
                        plr.teleport(new Location(plr.getWorld(), (double)warpMap.get("x"), (double)warpMap.get("y"), (double)warpMap.get("z"), (float)(double)warpMap.get("pitch"), (float)(double)warpMap.get("yaw")));
                        sender.sendMessage(ChatColor.GOLD+"Warping to "+ChatColor.RED+args[0]+ChatColor.GOLD+"...");
                    } catch (ClassCastException exc) {
                        log.log(Level.WARNING, "Error warping to "+args[0], exc);
                        sender.sendMessage(ChatColor.RED+"Error: "+ChatColor.DARK_RED+"That warp is corrupt. Please inform the admins.");
                    }
                }
            }
        }
        return true;
    }
    
    @Override
    public List<String> onTabComplete(CommandSender sender, Command cmd, String alias, String[] args) {
        List<String> result = new ArrayList<>();
        if (sender instanceof LivingEntity) {
            for (String s: new File(plugin.getDataFolder(), ((LivingEntity) sender).getWorld().getName()).list()) {
                if (args.length>0 && s.startsWith(args[0])) {
                    result.add(s);
                }
            }
        }
        return result;
    }
}
