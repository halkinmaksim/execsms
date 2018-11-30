ln -s /etc/init.d/execsms /etc/rc0.d/execsms
sudo update-rc.d execsms defaults 90



sudo rm -f execsms remove


*/1 * * * * /home/rxhf/execsms