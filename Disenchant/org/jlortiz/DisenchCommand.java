/*
 * Copyright (C) 2020 jlortiz
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

import java.util.Map.Entry;
import org.bukkit.ChatColor;
import org.bukkit.Material;
import org.bukkit.command.Command;
import org.bukkit.command.CommandExecutor;
import org.bukkit.command.CommandSender;
import org.bukkit.enchantments.Enchantment;
import org.bukkit.entity.Player;
import org.bukkit.inventory.ItemStack;
import org.bukkit.inventory.meta.EnchantmentStorageMeta;

/**
 *
 * @author jlortiz
 */
public class DisenchCommand implements CommandExecutor {
    @Override
    public boolean onCommand(CommandSender sender, Command cmd, String alias, String[] args) {
        if (!(sender instanceof Player)) {
            sender.sendMessage(ChatColor.RED+"This command is for players only!");
            return true;
        }
        Player p = (Player)sender;
        ItemStack held = p.getInventory().getItemInMainHand();
        if (held.getAmount() == 0) {
            sender.sendMessage(ChatColor.RED+"Try holding something.");
            return true;
        }
        if (held.getType() == Material.ENCHANTED_BOOK) {
            sender.sendMessage(ChatColor.RED+"Oh come on.");
            return true;
        }
        if (held.getEnchantments().isEmpty()) {
            sender.sendMessage(ChatColor.RED+"This item has no enchantments!");
            return true;
        }
        if (p.getInventory().first(Material.BOOK) == -1) {
            sender.sendMessage(ChatColor.RED+"You need to have a regular book to transfer enchantments to.");
            return true;
        }
        ItemStack book = p.getInventory().getItem(p.getInventory().first(Material.BOOK));
        if (p.getInventory().firstEmpty() == -1 && book.getAmount() > 1) {
            sender.sendMessage(ChatColor.RED+"Your inventory is full!");
            return true;
        }
        ItemStack enchBook = new ItemStack(Material.ENCHANTED_BOOK);
        EnchantmentStorageMeta meta = (EnchantmentStorageMeta)enchBook.getItemMeta();
        for (Entry<Enchantment,Integer> m: held.getEnchantments().entrySet()) {
            meta.addStoredEnchant(m.getKey(), m.getValue(), false);
            held.removeEnchantment(m.getKey());
        }
        if (book.getAmount() == 1) {
            p.getInventory().removeItem(book);
        }
        book.setAmount(book.getAmount() - 1);
        enchBook.setItemMeta(meta);
        p.getInventory().addItem(enchBook);
        sender.sendMessage(ChatColor.DARK_GREEN+"Enchantments have been transferred to book.");
        return true;
    }
}
