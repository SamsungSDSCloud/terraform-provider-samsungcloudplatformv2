provider "samsungcloudplatformv2" {
}

resource "samsungcloudplatformv2_filestorage_volume" "vol" {                                                                                     
    name      = var.volume_name                                                                                                                    
    protocol  = var.protocol                                                                                                                       
    type_name = var.type_name      
    cifs_password = var.cifs_password
    file_unit_recovery_enabled = var.file_unit_recovery_enabled
    tags = var.tags                                                                                                                
    zone = var.zone                                                                                                                           
  }                                                                                                                                                
                                                                                                                                                   
  resource "samsungcloudplatformv2_filestorage_access_rule" "node1" {                                                                              
    file_storage_id = samsungcloudplatformv2_filestorage_volume.vol.id                                                                             
    object_type     = var.object_type                                                                                                              
    object_id       = var.object_id                                                                                                                
  }